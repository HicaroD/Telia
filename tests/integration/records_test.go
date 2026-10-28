package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/tests/compiler"
)

func TestHeaderOwnedRecordLiteralUsesCLayoutAndZeroesOmittedFields(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureRecord"]
struct Record {
value i32
omitted i32
}
@[link_name="fixture_record_check"]
fn check(record *Record) i32
}
fn main() i32 {
record := Record.{value: 42}
if check(&record) != 1 {
return 1
}
return 0
}
`)
}

func runRecordFixture(t *testing.T, source string) {
	t.Helper()
	exe, diags := compileRecordFixture(t, source)
	if len(diags.Diags) != 0 {
		t.Fatalf("header-owned record: %v", diags.Diags)
	}
	defer removeTestBinary(t, exe)
	if output, err := compiler.RunBinary(exe); err != nil {
		t.Fatalf("record C fixture: %v: %s", err, output)
	}
}

func removeTestBinary(t *testing.T, path string) {
	t.Helper()
	if err := os.RemoveAll(filepath.Dir(path)); err != nil {
		t.Errorf("remove test binary: %v", err)
	}
}

func compileRecordFixture(
	t *testing.T,
	source string,
	cflags ...string,
) (string, *diagnostics.Collector) {
	t.Helper()
	header, err := filepath.Abs("testdata/records.h")
	if err != nil {
		t.Fatal(err)
	}
	fixture, err := filepath.Abs("testdata/records.c")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "main.t")
	if strings.Contains(source, "%q") {
		source = fmt.Sprintf(source, header)
	}
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	flags := append(
		[]string{
			fixture,
			"-Werror=incompatible-pointer-types",
			"-Werror=implicit-function-declaration",
		},
		cflags...)
	return compiler.CompileOnlyWithCFlags(path, flags...)
}

func TestHeaderOwnedRecordsRoundTripByValue(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureRecord"]
struct Record {
value i32
omitted i32
}
@[link_name="fixture_record_roundtrip"]
fn roundtrip(record Record) Record
}
fn main() i32 {
record := roundtrip(Record.{value: 40})
if record.value != 42 or record.omitted != 7 {
return 1
}
return 0
}
`)
}

func TestHeaderOwnedUnionLiteralSelectsMember(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureUnion"]
union Value {
integer i32
real f64
}
@[link_name="fixture_union_check"]
fn check(value Value) i32
}
fn main() i32 {
value := Value.{integer: 42}
if check(value) != 1 or value.integer != 42 {
return 1
}
return 0
}
`)
}
