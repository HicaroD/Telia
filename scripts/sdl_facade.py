"""Generate conservative Telia adapters; never infer ownership from char pointers."""
import re


# These calls consume/copy text before returning. Adding an entry requires an
# upstream ownership review, not just a matching C signature.
TEXT_CALLS = {
    "SDL_CreateWindow": "handle",
    "SDL_SetWindowTitle": "bool",
    "SDL_RenderDebugText": "bool",
    "SDL_SetClipboardText": "bool",
    "SDL_OpenURL": "bool",
    "SDL_LoadBMP": "handle",
    "SDL_SaveBMP": "bool",
    "SDL_GetWindowTitle": "borrowed",
    "SDL_GetError": "borrowed",
    "SDL_GetCurrentVideoDriver": "nullable_borrowed",
    "SDL_GetVideoDriver": "nullable_borrowed",
    "SDL_GetRendererName": "nullable_borrowed",
    "SDL_GetKeyName": "borrowed",
    "SDL_GetScancodeName": "borrowed",
    "SDL_GetClipboardText": "owned",
    "SDL_GetPrimarySelectionText": "owned",
    "SDL_SetPrimarySelectionText": "bool",
    "SDL_SetHint": "bool",
    "SDL_SetHintWithPriority": "bool",
    "SDL_ResetHint": "bool",
    "SDL_GetHint": "nullable_borrowed",
    "SDL_GetPlatform": "borrowed",
    "SDL_GetRevision": "borrowed",
    "SDL_GetAudioFormatName": "borrowed",
    "SDL_GetPixelFormatName": "borrowed",
}

# SDL owns these returned name/path views. A null result must remain distinct
# from a present empty value, and the facade must neither copy nor free it.
TEXT_CALLS.update(dict.fromkeys((
    "SDL_GetThreadName", "SDL_GetAudioDriver", "SDL_GetCurrentAudioDriver",
    "SDL_GetAudioDeviceName", "SDL_GetCameraDriver", "SDL_GetCurrentCameraDriver",
    "SDL_GetCameraName", "SDL_GetDisplayName", "SDL_GetRenderDriver",
    "SDL_GetSensorNameForID", "SDL_GetSensorName", "SDL_GetJoystickNameForID",
    "SDL_GetJoystickPathForID", "SDL_GetJoystickName", "SDL_GetJoystickPath",
    "SDL_GetJoystickSerial", "SDL_GetGamepadNameForID", "SDL_GetGamepadPathForID",
    "SDL_GetGamepadName", "SDL_GetGamepadPath", "SDL_GetGamepadSerial",
    "SDL_GetGamepadStringForType", "SDL_GetGamepadStringForAxis",
    "SDL_GetGamepadStringForButton", "SDL_GetGamepadAppleSFSymbolsNameForButton",
    "SDL_GetGamepadAppleSFSymbolsNameForAxis", "SDL_GetKeyboardNameForID",
    "SDL_GetMouseNameForID", "SDL_GetTouchDeviceName", "SDL_GetBasePath",
    "SDL_GetUserFolder", "SDL_GetGPUDriver", "SDL_GetGPUDeviceDriver",
    "SDL_GetHapticNameForID", "SDL_GetHapticName",
), "nullable_borrowed"))

MANUAL_CALLS = {
    "SDL_CreateRenderer": "create_renderer selects the default; create_renderer_named accepts an explicit name",
    "SDL_RenderDebugTextFormat": "render_debug_i32 is a typed alternative; arbitrary C formatting stays raw-only",
    "SDL_CreateThread": "create_thread bridges a cc=c callback with native scalar types and releases its trampoline state",
    "SDL_WaitThread": "wait_thread returns a native i32 exit status",
    "SDL_GetWindowSize": "get_window_size returns a native Size value instead of C output slots",
    "SDL_IOFromMem": "io_from_bytes borrows nonempty writable u8[] storage until the stream is closed",
    "SDL_ReadIO": "read_bytes bounds the read by a u8[] length and accepts empty slices",
}

SCALARS = {
    "c::c_int": "i32", "c::c_uint": "u32",
    "c::c_short": "i16", "c::c_ushort": "u16",
    "c::c_schar": "i8", "c::c_uchar": "u8",
    "c::c_longlong": "i64", "c::c_ulonglong": "u64",
    "c::c_float": "f32", "c::c_double": "f64",
    "c::c_size_t": "u64", "c::c_ptrdiff_t": "i64",
    "c::c_intptr_t": "i64", "c::c_uintptr_t": "u64",
}


def native_type(ty):
    if "c::" not in ty:
        return ty
    if ty in SCALARS:
        return SCALARS[ty]
    raise ValueError("C character, target-dependent scalar, or C scalar storage requires a deliberate adapter")


def callback_type(ty, generator, seen=None):
    ty = ty.removeprefix("const ").lstrip("*")
    seen = set() if seen is None else seen
    if ty in seen:
        return False
    seen.add(ty)
    declaration = generator.callbacks.get(ty, "")
    if "fn(" in declaration:
        return True
    alias = re.search(r"type \w+ = (.+)", declaration)
    return bool(alias and callback_type(alias[1], generator, seen))


def generate(generator, nodes, constants, text_constants):
    types = []
    for text in list(generator.declarations.values()) + list(generator.callbacks.values()):
        for name in re.findall(r"^(?:  )?(?:opaque|struct|union|enum|type) (\w+)", text, re.M):
            types.append("type " + name + " = raw::" + name)
    text = 'package sdl\n\nuse "vendor::sdl::raw"\nuse "std::strings"\n\n'
    text += "\n".join(types) + "\n\n"
    # Keep constants header-owned, with native scalar spelling. Text macros need
    # separate adapters: a C pointer must never masquerade as a Telia string.
    public_constants = []
    constant_inventory = []
    native_text_constants = []
    for declaration in constants:
        match = re.search(r"const (\w+) (.+)", declaration)
        name, ty = match.groups()
        if name in text_constants:
            native_text_constants.append("const " + name + " string = " + text_constants[name])
            constant_inventory.append({"name": name, "status": "native", "policy": "static_text"})
            continue
        try:
            native = native_type(ty)
            public_constants.append(declaration.replace(" " + ty + "\n", " " + native + "\n"))
            constant_inventory.append({"name": name, "status": "native"})
        except ValueError as error:
            constant_inventory.append({"name": name, "status": "raw_only", "reason": str(error)})
    text += '@[header="SDL3/SDL.h"]\nextern {\n' + "\n".join(public_constants) + "}\n\n"
    text += "\n".join(native_text_constants) + "\n\n"
    inventory = []
    for node in nodes:
        if node["kind"] != "FunctionDecl" or not node.get("name", "").startswith("SDL_"):
            continue
        name = node["name"]
        if name in MANUAL_CALLS:
            inventory.append({"name": name, "status": "adapted", "policy": MANUAL_CALLS[name]})
            continue
        try:
            declaration = generator.function(node)
            signature = re.search(r"fn (\w+)\((.*)\)(?: (.+))?\n", declaration)
            public, params, result = signature.groups()
            result = result or ""
            params = [p.split(" ", 1) for p in params.split(", ")] if params else []
            mode = TEXT_CALLS.get(name)
            if node.get("variadic"):
                raise ValueError("C varargs require a typed alternative")
            if any("c::c_va_list" in ty for _, ty in params):
                raise ValueError("C variadic cursor requires a typed alternative")
            if any("CString" in ty for _, ty in params) or "CString" in result:
                if not mode:
                    raise ValueError("text/buffer ownership, nullability, or retention requires a reviewed adapter")
            # Callback registration needs an explicit lifetime and calling-convention contract.
            if any(callback_type(ty, generator) for _, ty in params):
                raise ValueError("callback lifetime and calling convention require a deliberate adapter")
            native_params, args, setup = [], [], []
            for param, ty in params:
                if ty == "c::ConstCString":
                    native_params.append(param + " string")
                    setup += [f"  {param}_c, {param}_err := strings::to_c({param})",
                              f"  if {param}_err != nil {{", f"    return {failure_value(mode)}{param}_err", "  }",
                              f"  defer strings::free_c(&{param}_c)"]
                    args.append(f"strings::c_view(&{param}_c)")
                else:
                    native_params.append(param + " " + native_type(ty))
                    args.append(param)
            call = "raw::" + public + "(" + ", ".join(args) + ")"
            if mode == "bool":
                returns = "error"
                body = ["  if " + call + " == false {", f'    return error("SDL {public} failed")', "  }", "  return nil"]
            elif mode == "handle":
                returns = "(" + native_type(result) + ", error)"
                body = ["  result := " + call, "  if result == nil {", f'    return nil, error("SDL {public} failed")', "  }", "  return result, nil"]
            elif mode in ("borrowed", "nullable_borrowed"):
                returns = "OptionalString" if mode == "nullable_borrowed" else "string"
                if setup:
                    returns = "(" + returns + ", error)"
                body = ["  result := " + call, "  if result == nil {"]
                body += [('    return OptionalString{value: "", present: false}' if mode == "nullable_borrowed" else '    return ""') + (", nil" if setup else ""), "  }"]
                value = "OptionalString{value: strings::from_c(result), present: true}" if mode == "nullable_borrowed" else "strings::from_c(result)"
                body += ["  return " + value + (", nil" if setup else "")]
            elif mode == "owned":
                returns = "(strings::Owned, error)"
                body = ["  result := " + call, "  if result == nil {", f'    return strings::Owned{{text: ""}}, error("SDL {public} failed")', "  }", "  defer raw::free_text(result)", "  return strings::clone(strings::from_c(result))"]
            else:
                returns = native_type(result)
                body = [("  return " if result else "  ") + call]
            text += "fn " + public + "(" + ", ".join(native_params) + ")" + (" " + returns if returns else "") + " {\n"
            text += "\n".join(setup + body) + "\n}\n\n"
            inventory.append({"name": name, "status": "native", "policy": mode or "forward"})
        except ValueError as error:
            inventory.append({"name": name, "status": "raw_only", "reason": str(error)})
    return text.rstrip() + "\n", inventory, constant_inventory


def failure_value(mode):
    if mode == "bool":
        return ""
    if mode == "handle":
        return "nil, "
    if mode == "nullable_borrowed":
        return 'OptionalString{value: "", present: false}, '
    if mode == "borrowed":
        return '"", '
    raise ValueError("missing text-conversion failure policy")


def type_inventory(generator):
    """Make header-layout and callback exposure explicit, including record fields."""
    inventory = []
    for name, declaration in generator.declarations.items():
        if not declaration:
            continue
        row = {
            "name": name.removeprefix("SDL_"),
            "status": "header_identity",
            "policy": "original layout and ownership; pointer/callback fields remain low-level",
        }
        fields = re.findall(r"^    (\w+) (.+)$", declaration, re.M)
        if fields:
            row["fields"] = [{"name": field, "raw_type": ty} for field, ty in fields]
        inventory.append(row)
    for name, declaration in generator.callbacks.items():
        inventory.append({
            "name": name,
            "status": "raw_callback" if callback_type(name, generator) else "abi_alias",
            "declaration": declaration.strip(),
            "policy": "ABI identity preserved; no lifetime extension or ownership conversion",
        })
    return inventory
