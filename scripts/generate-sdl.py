#!/usr/bin/env python3
"""Generate header-owned Telia declarations from the pinned SDL3 C AST.

Clang is needed only for regeneration, never for ordinary Telia builds.
Unsupported declarations are recorded in coverage.json, not silently discarded.
"""
import argparse
import json
from pathlib import Path
import re
import subprocess
from sdl_facade import generate as generate_facade, type_inventory

ROOT = Path(__file__).resolve().parent.parent
SCALARS = {
    "void": "", "bool": "bool", "_Bool": "bool",
    "char": "c::c_char", "signed char": "c::c_schar", "unsigned char": "c::c_uchar",
    "short": "c::c_short", "unsigned short": "c::c_ushort",
    "int": "c::c_int", "unsigned int": "c::c_uint",
    "long": "c::c_long", "unsigned long": "c::c_ulong",
    "long long": "c::c_longlong", "unsigned long long": "c::c_ulonglong",
    "float": "c::c_float", "double": "c::c_double",
    "size_t": "c::c_size_t", "ptrdiff_t": "c::c_ptrdiff_t",
    "Sint8": "i8", "Uint8": "u8", "Sint16": "i16", "Uint16": "u16",
    "Sint32": "i32", "Uint32": "u32", "Sint64": "i64", "Uint64": "u64",
    "va_list": "c::c_va_list",
    "wchar_t": "c::c_wchar_t", "intptr_t": "c::c_intptr_t", "uintptr_t": "c::c_uintptr_t",
}
KEYWORDS = set("fn for while return extern package if elif else not and or type use defer struct union const enum match yield opaque true false nil rawptr bool int i8 i16 i32 i64 i128 uint u8 u16 u32 u64 u128 float f32 f64 string error".split())
CONSTANT_PREFIXES = tuple("SDL_" + prefix for prefix in (
    "INIT_", "WINDOW_", "WINDOWPOS_", "HINT_", "PROP_", "GPU_", "AUDIO_",
    "KMOD_", "HAPTIC_", "HAT_", "BUTTON_", "BLENDMODE_", "SURFACE_",
    "TRAYENTRY_", "PEN_", "MESSAGEBOX_", "ICONV_", "ALPHA_", "JOYSTICK_",
    "TOUCH_", "MOUSE_", "PI_", "FLT_", "STANDARD_", "RENDER_", "RENDERER_",
    "SECONDS_", "MS_", "US_", "NS_", "MAX_", "MIN_", "MAJOR_", "MINOR_", "MICRO_",
    "LIL_", "BIG_", "BYTEORDER", "VERSION", "CACHELINE_", "GL_", "GLOB_",
)) + ("SDLK_",)
EXTERNAL_TYPES = {"VkInstance", "VkPhysicalDevice", "VkSurfaceKHR", "VkAllocationCallbacks"}


def snake(name):
    if name == "SDL_log":
        return "log_f64"
    name = name.removeprefix("SDL_")
    name = re.sub(r"([A-Z]+)([A-Z][a-z])", r"\1_\2", name)
    return re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", name).lower()


def typename(name):
    return name.removeprefix("SDL_")


class Generator:
    def __init__(self, nodes):
        self.nodes = {n["name"]: n for n in nodes if (n.get("name", "").startswith("SDL_") or n.get("name") in EXTERNAL_TYPES) and n["kind"] == "TypedefDecl"}
        self.records = {n["name"]: n for n in nodes if n.get("name", "").startswith("SDL_") and n["kind"] in ("RecordDecl", "EnumDecl")}
        self.declarations = {}
        self.callbacks = {}
        self.pending = set()
        self.cnames = {}

    def ctype(self, text):
        if text == "XEvent":
            self.declarations[text] = '  @[link_name="XEvent"]\n  opaque XEvent\n'
            return text
        tag = re.fullmatch(r"(struct|union) ((?:SDL_|Vk)\w+)", text.strip())
        if tag and tag[2] not in self.nodes:
            self.declarations[tag[2]] = '  @[link_name="' + tag[1] + ' ' + tag[2] + '"]\n  opaque ' + typename(tag[2]) + '\n'
            return typename(tag[2])
        text = text.strip()
        if "(*)" in text:
            raise ValueError("anonymous callback requires a named callback declaration")
        array = re.fullmatch(r"(.+)\[(\d+)\]", text)
        if array:
            return self.ctype(array[1]) + "[" + array[2] + "]"
        if text.endswith("*const") or text.endswith("* const"):
            # Top-level parameter const does not change the function ABI.
            text = text[:text.rfind("const")].strip()
        if text.endswith("*"):
            pointee = text[:-1].strip()
            qualified = pointee.endswith("const") or ("*" not in pointee and pointee.startswith("const "))
            if pointee.endswith("const"):
                pointee = pointee[:-5].strip()
            if "*" not in pointee:
                pointee = re.sub(r"\bconst\b", "", pointee).strip()
            if pointee == "char":
                return "c::ConstCString" if qualified else "c::CString"
            if pointee == "void":
                if qualified:
                    self.declarations["void"] = '  @[link_name="void"]\n  opaque Void\n'
                    return "const *Void"
                return "rawptr"
            return ("const " if qualified else "") + "*" + self.ctype(pointee)
        if text in SCALARS:
            return SCALARS[text]
        text = re.sub(r"^(struct|union|enum) ", "", text)
        if text not in self.nodes:
            raise ValueError("unsupported C type: " + text)
        self.declare(text)
        return typename(text)

    def callback(self, name, text):
        match = re.fullmatch(r"(.+?)\s*\(\*\)\((.*)\)", text)
        if not match:
            raise ValueError("unsupported callback: " + text)
        args = [] if match[2] == "void" else match[2].split(", ")
        params = ", ".join("arg" + str(i) + " " + self.ctype(arg) for i, arg in enumerate(args))
        result = self.ctype(match[1])
        return '@[cc_macro="SDLCALL"]\ntype ' + name + " = fn(" + params + ")" + (" " + result if result else "") + "\n"

    def declare(self, name):
        if name in self.declarations or name in self.pending:
            return
        self.pending.add(name)
        try:
            node = self.nodes[name]
            text = node["type"]["qualType"]
            public = typename(name)
            if name in self.records:
                record = self.records[name]
                kind = record["kind"]
                cname = self.cnames.get(name, name)
                prefix = '  @[link_name="' + cname + '"]\n'
                if kind == "EnumDecl":
                    members = [n for n in record.get("inner", []) if n["kind"] == "EnumConstantDecl"]
                    common = __import__("os").path.commonprefix([n["name"] for n in members])
                    common = common[:common.rfind("_") + 1]
                    fields = []
                    for member in members:
                        membername = member["name"][len(common):]
                        if not membername or membername[0].isdigit():
                            membername = member["name"].removeprefix("SDL_")
                        fields.append('    @[link_name="' + member["name"] + '"]\n    ' + membername)
                    declaration = prefix + "  enum " + public + " {\n" + "\n".join(fields) + "\n  }\n"
                elif not record.get("completeDefinition"):
                    declaration = prefix + "  opaque " + public + "\n"
                else:
                    fields = []
                    anonymous = iter(n for n in record.get("inner", []) if n["kind"] == "RecordDecl")
                    for field in record.get("inner", []):
                        if field["kind"] != "FieldDecl":
                            continue
                        fieldname = field.get("name")
                        if not fieldname:
                            raise ValueError("anonymous record member in " + name)
                        fieldtype = field["type"]["qualType"]
                        if "(unnamed" in fieldtype or "(anonymous" in fieldtype:
                            nested = next(anonymous)
                            nestedname = name + fieldname[0].upper() + fieldname[1:]
                            self.nodes[nestedname] = {"type": {"qualType": nested["tagUsed"] + " " + nestedname}}
                            self.records[nestedname] = dict(nested, name=nestedname)
                            self.cnames[nestedname] = "__typeof__(((" + cname + " *)0)->" + fieldname + ")"
                            self.declare(nestedname)
                            fieldtype = typename(nestedname)
                        elif "(*)" in fieldtype:
                            cbname = public + "_" + fieldname + "Callback"
                            self.callbacks[cbname] = self.callback(cbname, fieldtype)
                            fieldtype = cbname
                        else:
                            fieldtype = self.ctype(fieldtype)
                        fields.append("    " + fieldname + " " + fieldtype)
                    declaration = prefix + "  " + record["tagUsed"] + " " + public + " {\n" + "\n".join(fields) + "\n  }\n"
            elif "(*)" in text:
                self.callbacks[public] = self.callback(public, text)
                declaration = ""
            else:
                self.callbacks[public] = "type " + public + " = " + self.ctype(text) + "\n"
                declaration = ""
            self.declarations[name] = declaration
        finally:
            self.pending.remove(name)

    def function(self, node):
        params = []
        for i, param in enumerate(n for n in node.get("inner", []) if n["kind"] == "ParmVarDecl"):
            name = param.get("name", "arg" + str(i))
            if name in KEYWORDS:
                name += "_"
            params.append(name + " " + self.ctype(param["type"]["qualType"]))
        if node.get("variadic"):
            params.append("args @c ...")
        result = self.ctype(node["type"]["qualType"].split("(")[0].strip())
        name = node["name"]
        return '  @[link_name="' + name + '"]\n  fn ' + snake(name) + "(" + ", ".join(params) + ")" + (" " + result if result else "") + "\n"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--include", type=Path, required=True, help="pinned artifact include directory")
    parser.add_argument("--output", type=Path, default=ROOT / "base/vendor/sdl")
    parser.add_argument("--clang", default="clang")
    args = parser.parse_args()
    header = args.include / "SDL3/SDL_version.h"
    if not header.exists() or not re.search(r"#define SDL_MICRO_VERSION\s+28\b", header.read_text()):
        parser.error("expected SDL 3.2.28 headers")
    command = [args.clang, "-x", "c", "-std=c11", "-DSDL_MAIN_HANDLED", "-I" + str(args.include)]
    unit = '#include <SDL3/SDL.h>\n#include <SDL3/SDL_main.h>\n#include <SDL3/SDL_vulkan.h>\n'
    result = subprocess.run(command + ["-Xclang", "-ast-dump=json", "-fsyntax-only", "-"], input=unit, text=True, capture_output=True, check=True)
    nodes = json.loads(result.stdout)["inner"]
    # These public entry points are macros, not FunctionDecls in the native AST.
    # Preserve the documented signatures and let SDL's header supply the hidden
    # platform C-runtime arguments instead of calling the Runtime APIs directly.
    for name, params in (
        ("SDL_CreateThread", (("entry", "SDL_ThreadFunction"), ("name", "const char *"), ("data", "void *"))),
        ("SDL_CreateThreadWithProperties", (("props", "SDL_PropertiesID"),)),
    ):
        nodes.append({
            "kind": "FunctionDecl", "name": name,
            "type": {"qualType": "SDL_Thread *(" + ", ".join(ty for _, ty in params) + ")"},
            "inner": [{"kind": "ParmVarDecl", "name": param, "type": {"qualType": ty}} for param, ty in params],
        })
    generator = Generator(nodes)
    functions = []
    inventory = []
    for node in nodes:
        if node["kind"] != "FunctionDecl" or not node.get("name", "").startswith("SDL_"):
            continue
        try:
            declaration = generator.function(node)
            functions.append(declaration)
            inventory.append({"name": node["name"], "status": "bound"})
        except ValueError as error:
            inventory.append({"name": node["name"], "status": "unsupported", "reason": str(error)})
    # Publish enum members as scalar constants too: SDL event payloads and flags
    # often use integer storage rather than the corresponding C enum type.
    constants = []
    text_constants = {}
    for record in generator.records.values():
        if record["kind"] != "EnumDecl":
            continue
        for member in record.get("inner", []):
            if member["kind"] == "EnumConstantDecl":
                name = member["name"]
                if name.startswith("SDL_"):
                    constants.append('  @[link_name="' + name + '"]\n  const ' + name.removeprefix("SDL_") + ' u32\n')
    macros = subprocess.run(command + ["-dM", "-E", "-"], input=unit, text=True, capture_output=True, check=True).stdout
    candidates = []
    for line in macros.splitlines():
        match = re.match(r"#define (\w+) (.+)", line)
        if match and match[1].startswith(CONSTANT_PREFIXES) and "_renamed_" not in match[2]:
            candidates.append(match[1])
    probes = unit + "\n".join("__auto_type telia_constant_" + name + " = " + name + ";" for name in sorted(candidates))
    result = subprocess.run(command + ["-Xclang", "-ast-dump=json", "-fsyntax-only", "-"], input=probes, text=True, capture_output=True)
    if result.returncode:
        parser.error("cannot inspect SDL constants:\n" + result.stderr)
    for node in json.loads(result.stdout)["inner"]:
        if node["kind"] != "VarDecl" or not node.get("name", "").startswith("telia_constant_"):
            continue
        name = node["name"].removeprefix("telia_constant_")
        texttype = node["type"].get("desugaredQualType", node["type"]["qualType"])
        if texttype.startswith("enum "):
            texttype = "unsigned int"
        ty = generator.ctype(texttype)
        ty = {
            "c::c_int": "i32", "c::c_uint": "u32",
            "c::c_short": "i16", "c::c_ushort": "u16",
            "c::c_schar": "i8", "c::c_uchar": "u8",
            "c::c_longlong": "i64", "c::c_ulonglong": "u64",
            "c::c_float": "f32", "c::c_double": "f64",
        }.get(ty, ty)
        if ty == "c::CString":
            ty = "c::ConstCString"
        if ty == "c::ConstCString":
            def literals(node):
                if node.get("kind") == "StringLiteral":
                    yield node["value"]
                for child in node.get("inner", []):
                    yield from literals(child)
            values = list(literals(node))
            if len(values) == 1:
                text_constants[name.removeprefix("SDL_")] = values[0]
        constants.append('  @[link_name="' + name + '"]\n  const ' + name.removeprefix("SDL_") + ' ' + ty + '\n')
    # SDL_main.h is intentionally included explicitly, with main handling disabled.
    text = 'package raw\n\nuse "std::ffi::c"\n\n'
    text += "\n".join(generator.callbacks.values())
    text += '\n@[header="SDL3/SDL.h"]\nextern {\n}\n\n@[header="SDL3/SDL_vulkan.h"]\nextern {\n}\n\n@[header="SDL3/SDL_main.h", default_cc_macro="SDLCALL"]\nextern {\n'
    text += "\n".join(generator.declarations.values()) + "\n" + "\n".join(constants) + "\n" + "\n".join(functions)
    text += '}\n\n@[default_cc="c"]\nextern {\n  @[link_name="telia_sdl_free_text"]\n  fn free_text(text c::CString)\n}\n'
    args.output.mkdir(parents=True, exist_ok=True)
    (args.output / "raw").mkdir(exist_ok=True)
    (args.output / "raw" / "sdl.t").write_text(text)
    facade, native_inventory, constant_inventory = generate_facade(generator, nodes, constants, text_constants)
    (args.output / "sdl.t").write_text(facade)
    (args.output / "coverage.json").write_text(json.dumps({"version": "3.2.28", "functions": inventory, "native_functions": native_inventory, "native_constants": constant_inventory, "native_types": type_inventory(generator)}, indent=2) + "\n")
    print(str(sum(item["status"] == "bound" for item in inventory)) + "/" + str(len(inventory)) + " functions generated; see coverage.json for gaps")


if __name__ == "__main__":
    main()
