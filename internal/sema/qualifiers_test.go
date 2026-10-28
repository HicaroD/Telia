package sema_test

import "testing"

func TestConstPointerCannotLoseItsQualifier(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn drop(pointer const *i32) *i32 {\nreturn pointer\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Const qualifier removed" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want const removal diagnostic, got %v", c.Diags)
	}
}

func TestConstArrayCannotBeBorrowedAsMutableSlice(t *testing.T) {
	c := checkDiagnostics(t, "package main\nstruct Record {\nvalues i32[2]\n}\nfn write(record const *Record) {\nvalues := record.values[:]\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Const pointee mutation" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want const slice diagnostic, got %v", c.Diags)
	}
}

func TestConstRecordArrayElementsCannotBeMutated(t *testing.T) {
	c := checkDiagnostics(t, "package main\nstruct Record {\nvalues i32[2]\n}\nfn write(record const *Record) {\nrecord.values[0] = 42\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Const pointee mutation" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want const array diagnostic, got %v", c.Diags)
	}
}

func TestConstRecordFieldCannotEscapeAsMutablePointer(t *testing.T) {
	c := checkDiagnostics(t, "package main\nstruct Record {\nvalue i32\n}\nfn write(record const *Record) {\npointer *i32 := &record.value\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Const pointee mutation" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want const escape diagnostic, got %v", c.Diags)
	}
}

func TestConstPointerRejectsWritingItsImmediatePointee(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn write(value const *i32) {\n*value = 42\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Const pointee mutation" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want const mutation diagnostic, got %v", c.Diags)
	}
}

func TestConstPointerRejectsMutatingRecordFields(t *testing.T) {
	c := checkDiagnostics(t, "package main\nstruct Record {\nvalue i32\n}\nfn write(record const *Record) {\nrecord.value = 42\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Const pointee mutation" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want const field diagnostic, got %v", c.Diags)
	}
}
