package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

// An opt-in native measurement, not a performance assertion. Compilation and
// key construction are outside the timer. Run with -v to retain the CSV rows.
func TestContainerBenchmarks(t *testing.T) {
	if os.Getenv("TELIA_CONTAINER_BENCH") != "1" {
		t.Skip("set TELIA_CONTAINER_BENCH=1 for native CSV measurements")
	}
	header, err := filepath.Abs("testdata/container_benchmark.h")
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"integer", "string", "record"} {
		for _, container := range []string{"maps", "sets"} {
			t.Run(container+"/"+kind, func(t *testing.T) {
				keyType, key := "u64", "index"
				if kind == "string" {
					keyType, key = "string", "strings::from_c(_bench_key(index))"
				} else if kind == "record" {
					keyType, key = "Record", "Record{index: index, label: \"record-key\"}"
				}
				create := container + "::create<" + keyType + ">()"
				insert := "sets::insert(&table, keys[index]) @fail"
				lookup := "if sets::contains(&table, keys[index]) {\n checksum = checksum + 1\n }"
				iterate := "for key in table {\n checksum = checksum + 1\n }"
				if container == "maps" {
					create = "maps::create<" + keyType + ", u64>()"
					insert = "maps::insert(&table, keys[index], index) @fail"
					lookup = "value, found := maps::get(&table, keys[index])\n if found {\n checksum = checksum + value\n }"
					iterate = "for key, value in table {\n checksum = checksum + value\n }"
				}
				loop := func(body string) string {
					return "index = 0\n while index < size {\n" + body + "\n index = index + 1\n }\n"
				}
				source := `package main
use "std::maps"
use "std::sets"
use "std::strings"
use "std::ffi::c"
struct Record {
 index u64
 label string
}
@[default_cc="c"]
extern {
 fn _bench_key(index u64) c::ConstCString
 fn _bench_start()
 fn _bench_stop(operation u64, count u64)
 fn _bench_report(size u64, checksum u64)
}
fn run(size u64) {
 keys KEYTYPE[8192] := .{}
 index u64 := 0
 while index < size {
  keys[index] = KEY
  index = index + 1
 }
 checksum u64 := 0
 repeat u64 := 0
 while repeat < 32 {
  table := CREATE
  CONTAINER::reserve(&table, size) @fail
  _bench_start()
  INSERT
  _bench_stop(0, size)
  _bench_start()
  LOOKUP
  _bench_stop(1, size)
  _bench_start()
  ITERATE
  _bench_stop(2, size)
  _bench_start()
  REMOVE
  _bench_stop(3, size)
  _bench_start()
  CHURN
  _bench_stop(4, size)
  CONTAINER::free(&table)
  table = CREATE
  _bench_start()
  INSERT
  _bench_stop(5, size)
  CONTAINER::free(&table)
  repeat = repeat + 1
 }
 _bench_report(size, checksum)
 return
}
fn main() {
 run(64)
 run(1024)
 run(8192)
 return
}
`
				source = strings.NewReplacer(
					"KEYTYPE", keyType, "KEY", key, "CREATE", create,
					"CONTAINER", container, "INSERT", loop(insert),
					"LOOKUP", loop(lookup), "ITERATE", iterate,
					"REMOVE", loop(container+"::remove(&table, keys[index])"),
					"CHURN", loop(insert+"\n"+container+"::remove(&table, keys[index])"),
				).Replace(source)
				file := filepath.Join(t.TempDir(), "main.t")
				if err := os.WriteFile(file, []byte(source), 0600); err != nil {
					t.Fatal(err)
				}
				path, diags := compiler.CompileOnlyWithCFlags(file, "-O3", "-include", header)
				if diags.HasErrors() {
					t.Fatalf("benchmark compile: %v", diags.Diags)
				}
				defer removeTestBinary(t, path)
				out, err := exec.Command(path).CombinedOutput()
				if err != nil {
					t.Fatalf("benchmark execution: %v: %s", err, out)
				}
				for _, row := range strings.Split(strings.TrimSpace(string(out)), "\n") {
					fmt.Printf("BENCH,%s,%s,%s\n", container, kind, row)
				}
			})
		}
	}
}
