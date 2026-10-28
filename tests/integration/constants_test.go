package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestNamedConstantIsUsableInAnExpression(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := "package main\nconst Answer i32 = 42\nfn main() i32 {\nif Answer != 42 {\nreturn 1\n}\nreturn 0\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("named constant: %v", diags.Diags)
	}
}

func TestConstantsEvaluateFlagExpressionsAtCompileTime(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := "package main\nconst Mask u64 = 8 | 32\nconst Selected u64 = Mask & ~8\nfn main() i32 {\nif Selected != 32 {\nreturn 1\n}\nreturn 0\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("compile-time flag expression: %v", diags.Diags)
	}
}

func TestLocalConstantsInferTypesAndRespectBlockScope(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := "package main\nconst Base = 40\nfn main() i32 {\nconst Answer = Base + 2\n{\nconst Answer = 3\nif Answer != 3 {\nreturn 1\n}\n}\nif Answer != 42 {\nreturn 2\n}\nreturn 0\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("block-scoped inferred constants: %v", diags.Diags)
	}
}

func TestImportedConstantsCanInitializeOtherConstants(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "flags"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "flags", "flags.t"), []byte("package flags\nconst Mask u64 = 8 | 32\n"), 0600); err != nil {
		t.Fatal(err)
	}
	source := "package main\nuse \"pkg::flags\"\nconst Selected u64 = flags::Mask & ~8\nfn main() i32 {\nif flags::Mask != 40 {\nreturn 1\n}\nif Selected != 32 {\nreturn 2\n}\nreturn 0\n}\n"
	if err := os.WriteFile(filepath.Join(dir, "main.t"), []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompilePackage(dir)
	if len(diags.Diags) != 0 {
		t.Fatalf("imported constants: %v", diags.Diags)
	}
}

func TestGenericFunctionsHaveIndependentLocalConstants(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := "package main\nconst Base = 40\nfn answer<T>() T {\nconst Value T = Base + 2\nreturn Value\n}\nfn main() i32 {\nsmall := answer<u8>()\nlarge := answer<u64>()\nif small != 42 {\nreturn 1\n}\nif large != 42 {\nreturn 2\n}\nreturn 0\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("specialized constants must not share mutable state: %v", diags.Diags)
	}
}

func TestIntegerConstantsPreserveTheirFullDeclaredWidth(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := "package main\nconst Full u128 = 340282366920938463463374607431768211455\nconst Low u64 = 18446744073709551615\nconst Minimum i128 = -170141183460469231731687303715884105728\nfn main() i32 {\nif Full != ~zero<u128>() {\nreturn 1\n}\nif Low != ~zero<u64>() {\nreturn 2\n}\nif Minimum >= 0 {\nreturn 3\n}\nreturn 0\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("full-width integer constants: %v", diags.Diags)
	}
}

func TestConstantsSupportBooleanStringAndFloatingValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := "package main\nconst Enabled = true\nconst Label = \"Telia\"\nconst Step f32 = 1.5\nconst Half f64 = 1.0 / 2.0\nfn main() i32 {\nif not Enabled {\nreturn 1\n}\nif Label != \"Telia\" {\nreturn 2\n}\nif Step != 1.5 {\nreturn 3\n}\nif Half != 0.5 {\nreturn 4\n}\nreturn 0\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("scalar constants: %v", diags.Diags)
	}
}

func TestForwardConstantsUseTheirOwnFileImports(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "settings"), 0700); err != nil {
		t.Fatal(err)
	}
	for name, source := range map[string]string{
		"a_main.t":            "package main\nfn main() i32 {\nif Answer != 42 {\nreturn 1\n}\nreturn 0\n}\n",
		"z_constants.t":       "package main\nuse \"pkg::settings\"\nconst Answer = settings::Value\n",
		"settings/settings.t": "package settings\nconst Value = 42\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	_, diags := compiler.CompilePackage(dir)
	if len(diags.Diags) != 0 {
		t.Fatalf("forward constant must resolve imports in its declaring file: %v", diags.Diags)
	}
}

func TestImportedInferredConstantsAdaptWithoutChangingTheirDeclaration(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "settings"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "settings", "settings.t"), []byte("package settings\nconst Count = 42\n"), 0600); err != nil {
		t.Fatal(err)
	}
	source := "package main\nuse \"pkg::settings\"\nconst Small u8 = settings::Count\nfn main() i32 {\nlarge u64 := settings::Count\nif large != 42 {\nreturn 1\n}\nif Small != 42 {\nreturn 2\n}\nif large & settings::Count != 42 {\nreturn 3\n}\nreturn 0\n}\n"
	if err := os.WriteFile(filepath.Join(dir, "main.t"), []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompilePackage(dir)
	if len(diags.Diags) != 0 {
		t.Fatalf("imported constant uses must keep independent contextual types: %v", diags.Diags)
	}
}

func TestTypedFloatingConstantsRoundBeforeDependentArithmetic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := "package main\nconst Rounded f32 = 16777217.0\nconst Difference f32 = Rounded - 16777216.0\nfn main() i32 {\nif Difference != 0.0 {\nreturn 1\n}\nreturn 0\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("typed f32 constants must use their representable value: %v", diags.Diags)
	}
}
