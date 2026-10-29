package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestHeaderOwnedKeywordFieldsRoundTripThroughC(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "event.h")
	if err := os.WriteFile(header, []byte("typedef struct Event { unsigned int type; } Event;\nstatic inline unsigned int event_type(Event *event) { return event->type; }\n"), 0600); err != nil {
		t.Fatal(err)
	}
	source := fmt.Sprintf(`package main
@[header=%q, default_cc="c"]
extern {
  struct Event {
    type u32
  }
  fn event_type(event *Event) u32
}
fn main() i32 {
  event := Event{type: 256}
  event.type = event.type + 1
  if event_type(&event) != 257 {
    return 1
  }
  return 0
}
`, header)
	path := filepath.Join(dir, "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("keyword-named C field must round-trip: %v", diags.Diags)
	}
}
