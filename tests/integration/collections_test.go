package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestArraysCopyInlineElementsAndZeroInitialize(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/arrays.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("arrays should compile: %v", diags.Diags)
	}
	if want := "3\n20\n0\n99\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestCollectionBoundsChecksInDebugAndRelease(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, tc := range []struct{ name, body, want string }{
			{"array", "values i32[1] := .{42}\nx := values[1]", "index out of bounds"},
			{"empty array", "values i32[0] := .{}\nx := values[0]", "index out of bounds"},
			{"slice", "values i32[1] := .{}\nview := values[:]\nx := view[1]", "index out of bounds"},
			{"empty slice", "view := zero<i32[]>()\nx := view[0]", "index out of bounds"},
			{"range", "values i32[1] := .{}\nx := values[0:2] @fail", "slice range out of bounds"},
		} {
			t.Run(opt+"/"+tc.name, func(t *testing.T) {
				file := filepath.Join(t.TempDir(), "bounds.t")
				if err := os.WriteFile(file, []byte("package main\nfn main() i32 {\n"+tc.body+"\nreturn 0\n}\n"), 0600); err != nil {
					t.Fatal(err)
				}
				path, diags := compiler.CompileOnlyWithCFlags(
					file,
					opt,
					"-fsanitize=undefined",
					"-fno-sanitize-recover=all",
				)
				if len(diags.Diags) != 0 {
					t.Fatalf("bounds program should compile: %v", diags.Diags)
				}
				defer removeTestBinary(t, path)
				output, err := compiler.RunBinary(path)
				if err == nil || output != "runtime panic: "+tc.want+"\n" {
					t.Fatalf("want bounds panic, got %q, %v", output, err)
				}
			})
		}
	}
}

func TestCollectionsWithNativeSanitizers(t *testing.T) {
	sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
	if sanitizers == "" {
		sanitizers = "undefined"
	}
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range []string{"arrays.t", "slices.t", "collection_iteration.t", "collection_memory.t", "collection_evaluation.t", "collection_composition.t", "collection_ownership.t", "list_aliasing.t", "list_lifecycle.t", "list_failure.t", "list_empty_struct.t"} {
		t.Run(fixture, func(t *testing.T) {
			path, diags := compiler.CompileOnlyWithCFlags(
				"testdata/"+fixture,
				"-include",
				header,
				"-fsanitize="+sanitizers,
				"-fno-sanitize-recover=all",
				"-fno-omit-frame-pointer",
			)
			if len(diags.Diags) != 0 {
				t.Fatalf("sanitized compilation: %v", diags.Diags)
			}
			defer removeTestBinary(t, path)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			output, err := exec.CommandContext(ctx, path).CombinedOutput()
			if err != nil {
				t.Fatalf("sanitized execution: %v: %s", err, output)
			}
		})
	}
}

func TestCollectionsComposeWithStructsArraysPointersAndNestedLists(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/collection_composition.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("composed collections should compile: %v", diags.Diags)
	}
	if want := "record\n2\n9\nnested\nrecord\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestListsSupportEmptyStructElementsWithoutZeroSizeAllocation(t *testing.T) {
	path, diags := compiler.CompileOnlyWithCFlags(
		"testdata/list_empty_struct.t",
		"-fsanitize=undefined",
		"-fno-sanitize-recover=all",
	)
	if len(diags.Diags) != 0 {
		t.Fatalf("empty element structs should compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	output, err := compiler.RunBinary(path)
	if err != nil || output != "1\n1\n" {
		t.Fatalf("want one stored nonzero-size value, got %q, %v", output, err)
	}
}

func TestCollectionElementsExposePointersForExplicitResourceCleanup(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/collection_ownership.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("element pointers should compile: %v", diags.Diags)
	}
	if want := "owned\n0\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestIndexAndSliceOperandsAreEvaluatedOnceInSourceOrder(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/collection_evaluation.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("collection operands should compile: %v", diags.Diags)
	}
	if want := "source\nindex\n42\nsource\nindex\n42\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestListFailuresPreserveStorageAndReleaseAllocations(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	path, diags := compiler.CompileOnlyWithCFlags(
		"testdata/list_failure.t",
		"-include",
		header,
		"-fsanitize=undefined",
		"-fno-sanitize-recover=all",
	)
	if len(diags.Diags) != 0 {
		t.Fatalf("list failure fixture should compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	output, err := compiler.RunBinary(path)
	if err != nil {
		t.Fatalf("execution: %v: %s", err, output)
	}
	if want := "reserve failed\nappend failed\nclone failed\noverflow rejected\n2\n2\narena failed\n7\nreleased\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestCheckedSliceAccessReturnsErrorsWithoutMutation(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/slice_access.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("checked slice access should compile: %v", diags.Diags)
	}
	if want := "42\nmissing\nunchanged\n42\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestListRemovalClearAndCloneKeepExplicitValueOwnership(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/list_lifecycle.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("list lifecycle should compile: %v", diags.Diags)
	}
	if want := "20\n10\n30\n1\n40\n0\n0\n1\n10\n20\n30\n40\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestListsAppendAliasedSlicesAcrossHeapAndArenaGrowth(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/list_aliasing.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("alias-safe list append should compile: %v", diags.Diags)
	}
	if want := "1\n2\n3\n2\n3\n1\n2\n3\n2\n3\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestTeliaListsGrowAndExposeBorrowedViews(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/lists.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("Telia lists should compile: %v", diags.Diags)
	}
	if want := "3\n10\n77\n30\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericLayoutPrimitivesSupportTypedHeapStorage(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/collection_memory.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("typed storage should compile: %v", diags.Diags)
	}
	if want := "4\n4\n42\n0\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericArraysAndSlicesIterateByValue(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/collection_iteration.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("generic collection iteration should compile: %v", diags.Diags)
	}
	if want := "6\n1\n0\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestSlicesBorrowWritableStorageAndCheckRanges(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/slices.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("slices should compile: %v", diags.Diags)
	}
	if want := "3\n2\n77\n0\nbad range\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}
