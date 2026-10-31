package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSDL3NativeHintDistinguishesAbsentAndEmptyText(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  missing := sdl::get_hint("TELIA_NATIVE_ABSENT_HINT") @fail
  if missing.present or missing.value != "" {
    return 1
  }
  sdl::set_hint("TELIA_NATIVE_EMPTY_HINT", "") @fail
  found := sdl::get_hint("TELIA_NATIVE_EMPTY_HINT") @fail
  if found.present == false or found.value != "" {
    return 2
  }
  sdl::reset_hint("TELIA_NATIVE_EMPTY_HINT") @fail
  return 0
}
`)
}

func TestSDL3NativeThreadRunsACCallbackWithNativeScalarTypes(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
@[cc="c"]
fn worker(data rawptr) i32 {
  return 42
}
fn main() i32 {
  thread := sdl::create_thread(worker, "native-worker", nil) @fail
  if sdl::wait_thread(thread) != 42 {
    return 1
  }
  return 0
}
`)
}

func TestSDL3NativeTextRejectsEmbeddedZeroWithoutChangingTheWindow(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    return 1
  }
  defer sdl::quit()
  window := sdl::create_window("unchanged", 64, 64, 0) @fail
  defer sdl::destroy_window(window)
  text string := "bad\0title"
  err := sdl::set_window_title(window, text)
  if err == nil or err.msg != "embedded zero byte" {
    return 2
  }
  if sdl::get_window_title(window) != "unchanged" {
    return 3
  }
  invalid, create_err := sdl::create_window(text, 64, 64, 0)
  if invalid != nil or create_err == nil {
    return 4
  }
  return 0
}
`)
}

func TestSDL3NativeWindowSizeReturnsNativeValues(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    return 1
  }
  defer sdl::quit()
  window := sdl::create_window("size", 64, 48, 0) @fail
  defer sdl::destroy_window(window)
  size := sdl::get_window_size(window) @fail
  if size.width != 64 or size.height != 48 {
    return 2
  }
  return 0
}
`)
}

func TestSDL3NativeByteIOPreservesZerosAndEmptySlices(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  bytes u8[3] := .{65, 0, 66}
  stream := sdl::io_from_bytes(bytes[:]) @fail
  defer sdl::close_io(stream)
  output u8[3] := .{}
  empty := output[0:0] @fail
  if sdl::read_bytes(stream, empty) != 0 {
    return 1
  }
  if sdl::read_bytes(stream, output[:]) != 3 {
    return 2
  }
  if output[0] != 65 or output[1] != 0 or output[2] != 66 {
    return 3
  }
  return 0
}
`)
}

func TestSDL3NativePlatformAndDisplayNamesUseNativeText(t *testing.T) {
	runSDLSource(t, `package main
use "std::strings"
use "vendor::sdl"
fn main() i32 {
  if strings::len(sdl::get_platform()) == 0 {
    return 1
  }
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    return 2
  }
  defer sdl::quit()
  display := sdl::get_primary_display()
  name := sdl::get_display_name(display)
  if name.present == false or strings::len(name.value) == 0 {
    return 3
  }
  return 0
}
`)
}

func TestSDL3NativeClipboardTextIsIndependentlyOwned(t *testing.T) {
	runSDLSource(t, `package main
use "std::strings"
use "vendor::sdl"
fn main() i32 {
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    return 1
  }
  defer sdl::quit()
  sdl::set_clipboard_text("original") @fail
  before := sdl::get_num_allocations()
  text := sdl::get_clipboard_text() @fail
  defer strings::free(&text)
  if sdl::get_num_allocations() != before {
    return 2
  }
  sdl::set_clipboard_text("replacement") @fail
  if strings::view(&text) != "original" {
    return 3
  }
  return 0
}
`)
}

func TestSDL3NativeTextReportsAllocationFailureWithoutCallingSDL(t *testing.T) {
	if os.Getenv("TELIA_TEST_SDL") != "1" {
		t.Skip("run SDL setup, then set TELIA_TEST_SDL=1")
	}
	header := filepath.Join(t.TempDir(), "fail.h")
	if err := os.WriteFile(header, []byte(`#include <stdlib.h>
static void *fail_alloc(size_t size) { (void)size; return NULL; }
#define _TELIA_MALLOC fail_alloc
`), 0600); err != nil {
		t.Fatal(err)
	}
	runSDLSourceWithCFlags(t, `package main
use "vendor::sdl"
fn main() i32 {
  window, err := sdl::create_window("allocation failure", 64, 64, 0)
  if window != nil or err == nil {
    return 1
  }
  if err.msg != "allocation failed" {
    return 2
  }
  if sdl::get_error() != "" {
    return 3
  }
  return 0
}
`, "-include", header)
}

func TestSDL3NativeTextCleansUpWhenASecondConversionFails(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	runSDLSourceWithCFlags(t, `package main
use "vendor::sdl"
@[default_cc="c"]
extern {
  fn _test_fail_after(count i64)
  fn _test_live_allocations() i64
}
fn main() i32 {
  _test_fail_after(1)
  err := sdl::set_hint("TELIA_NATIVE_CLEANUP_HINT", "value")
  if err == nil or err.msg != "allocation failed" {
    return 1
  }
  if _test_live_allocations() != 0 {
    return 2
  }
  _test_fail_after(-1)
  missing := sdl::get_hint("TELIA_NATIVE_CLEANUP_HINT") @fail
  if missing.present {
    return 3
  }
  sdl::set_hint("TELIA_NATIVE_CLEANUP_HINT", "value") @fail
  if _test_live_allocations() != 0 {
    return 4
  }
  sdl::reset_hint("TELIA_NATIVE_CLEANUP_HINT") @fail
  return 0
}
`, "-include", header)
}
