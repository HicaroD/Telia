"""Boundary tests for generated native-vs-raw SDL API policy."""
import importlib.util
from pathlib import Path
import unittest

from sdl_facade import generate, type_inventory


SPEC = importlib.util.spec_from_file_location("generate_sdl", Path(__file__).with_name("generate-sdl.py"))
BINDINGS = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(BINDINGS)


def function(name, result="bool", params=(), variadic=False):
    return {
        "kind": "FunctionDecl", "name": name,
        "type": {"qualType": result + "(" + ", ".join(ty for _, ty in params) + ")"},
        "inner": [{"kind": "ParmVarDecl", "name": name, "type": {"qualType": ty}} for name, ty in params],
        "variadic": variadic,
    }


def facade(nodes, constants=(), text_constants=None):
    generator = BINDINGS.Generator(nodes)
    for node in nodes:
        if node["kind"] == "FunctionDecl":
            generator.function(node)
    return generate(generator, nodes, constants, text_constants or {})


class NativeFacadePolicyTests(unittest.TestCase):
    def test_record_fields_are_audited_without_replacing_the_header_layout(self):
        nodes = [
            {"kind": "RecordDecl", "name": "SDL_Record", "tagUsed": "struct", "completeDefinition": True,
             "inner": [{"kind": "FieldDecl", "name": "title", "type": {"qualType": "const char *"}}]},
            {"kind": "TypedefDecl", "name": "SDL_Record", "type": {"qualType": "struct SDL_Record"}},
        ]
        generator = BINDINGS.Generator(nodes)
        generator.declare("SDL_Record")
        text, _, _ = generate(generator, nodes, [], {})
        self.assertIn("type Record = raw::Record", text)
        self.assertNotIn("struct Record", text)
        inventory = type_inventory(generator)
        self.assertEqual(inventory[0]["fields"], [{"name": "title", "raw_type": "c::ConstCString"}])
        self.assertEqual(inventory[0]["status"], "header_identity")

    def test_aliases_cannot_hide_callback_registration_from_the_audit(self):
        nodes = [
            {"kind": "TypedefDecl", "name": "SDL_Notify", "type": {"qualType": "void (*)(void *)"}},
            {"kind": "TypedefDecl", "name": "SDL_NotifyAlias", "type": {"qualType": "SDL_Notify"}},
            function("SDL_Register", params=(("callback", "SDL_NotifyAlias"),)),
        ]
        text, coverage, _ = facade(nodes)
        self.assertEqual(coverage[0]["status"], "raw_only")
        self.assertNotIn("fn register", text)

    def test_unreviewed_text_and_retained_names_are_not_temporary_conversions(self):
        nodes = [function(name, params=(("name", "const char *"),)) for name in ("SDL_UnknownText", "SDL_SetScancodeName")]
        text, coverage, _ = facade(nodes)
        self.assertEqual([row["status"] for row in coverage], ["raw_only", "raw_only"])
        self.assertNotIn("fn unknown_text", text)
        self.assertNotIn("fn set_scancode_name", text)
        self.assertTrue(all(row.get("reason") for row in coverage))

    def test_binary_char_buffers_are_not_presented_as_strings(self):
        text, coverage, _ = facade([function("SDL_Buffer", params=(("data", "char *"),))])
        self.assertEqual(coverage[0]["status"], "raw_only")
        self.assertNotIn("fn buffer", text)

    def test_varargs_and_target_dependent_scalars_remain_explicitly_raw(self):
        nodes = [function("SDL_Print", params=(("format", "const char *"),), variadic=True), function("SDL_Long", result="long")]
        _, coverage, _ = facade(nodes)
        self.assertEqual([row["status"] for row in coverage], ["raw_only", "raw_only"])

    def test_scalar_values_use_native_types_but_pointer_storage_is_not_retyped(self):
        nodes = [function("SDL_Value", result="int", params=(("value", "float"),)), function("SDL_Output", params=(("value", "int *"),))]
        text, coverage, _ = facade(nodes)
        self.assertIn("fn value(value f32) i32", text)
        self.assertEqual([row["status"] for row in coverage], ["native", "raw_only"])

    def test_reviewed_text_frees_each_conversion_and_propagates_conversion_failure(self):
        text, coverage, _ = facade([function("SDL_SetHint", params=(("name", "const char *"), ("value", "const char *")))])
        self.assertEqual(coverage[0]["policy"], "bool")
        self.assertIn("fn set_hint(name string, value string) error", text)
        self.assertIn("defer strings::free_c(&name_c)", text)
        self.assertIn("defer strings::free_c(&value_c)", text)
        self.assertIn("return value_err", text)

    def test_text_constants_are_native_static_strings_not_pointer_aliases(self):
        declaration = '  @[link_name="SDL_HINT_RENDER_DRIVER"]\n  const HINT_RENDER_DRIVER c::ConstCString\n'
        text, _, constants = facade([], [declaration], {"HINT_RENDER_DRIVER": '"SDL_RENDER_DRIVER"'})
        self.assertIn('const HINT_RENDER_DRIVER string = "SDL_RENDER_DRIVER"', text)
        self.assertEqual(constants[0]["status"], "native")


if __name__ == "__main__":
    unittest.main()
