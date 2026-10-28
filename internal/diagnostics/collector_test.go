package diagnostics_test

import (
	"io"
	"os"
	"testing"

	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/source"
)

func TestCollectionDoesNotPrint(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	previous := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = previous }()
	c := diagnostics.New()
	c.ReportAndSave(diagnostics.Diag{Message: "bad value"})
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("collection printed %q", got)
	}
	if len(c.Diags) != 1 {
		t.Fatalf("want one diagnostic, got %d", len(c.Diags))
	}
}

func TestLegacyStageMessagesDoNotDuplicateLocationsWithSpaces(t *testing.T) {
	c := diagnostics.New()
	f := source.New("my file.t", []byte("bad"))
	c.ReportAndSave(
		diagnostics.Diag{
			Message: "my file.t:1:1: unknown value",
			Span:    source.Span{File: f, Start: 0, End: 3},
		},
	)
	if c.Diags[0].Message != "unknown value" {
		t.Fatalf("location remained in message: %q", c.Diags[0].Message)
	}
}
