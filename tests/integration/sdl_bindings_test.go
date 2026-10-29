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

func runSDLSource(t *testing.T, source string) {
	t.Helper()
	if os.Getenv("TELIA_TEST_SDL") != "1" {
		t.Skip("run SDL setup, then set TELIA_TEST_SDL=1")
	}
	root, err := filepath.Abs("../../base/vendor")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("T_VENDOR_PATH", root)
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	binary, diags := compiler.CompileOnly(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("compile SDL program: %v", diags.Diags)
	}
	t.Cleanup(func() { removeTestBinary(t, binary) })
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Env = append(
		os.Environ(),
		"SDL_VIDEO_DRIVER=dummy",
		"SDL_RENDER_DRIVER=software",
		"SDL_AUDIO_DRIVER=dummy",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("SDL program: %v\n%s", err, out)
	}
}

func TestSDL3KeyboardEventPayloadRoundTrips(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    return 1
  }
  defer sdl::quit()
  event := sdl::Event{type: 0}
  while sdl::poll_event(&event) {
  }
  key := sdl::KeyboardEvent{type: sdl::EventType::KEY_DOWN, scancode: sdl::Scancode::RIGHT, down: true, repeat: false}
  sent := sdl::Event{key: key}
  if sdl::push_event(&sent) == false {
    return 2
  }
  if sdl::poll_event(&event) == false {
    return 3
  }
  if event.type != sdl::EVENT_KEY_DOWN {
    return 4
  }
  if event.key.scancode != sdl::Scancode::RIGHT {
    return 5
  }
  if event.key.down == false or event.key.repeat {
    return 6
  }
  if sdl::poll_event(&event) {
    return 7
  }
  return 0
}
`)
}

func TestSDL3WideTextUsesNativeWcharStorage(t *testing.T) {
	runSDLSource(t, `package main
use "std::ffi::c"
use "vendor::sdl"
fn main() i32 {
  text c::c_wchar_t[3] := .{65, 66, 0}
  if sdl::wcslen(&text[0]) != 2 {
    return 1
  }
  copy := sdl::wcsdup(&text[0])
  if copy == nil {
    return 2
  }
  defer sdl::free(copy)
  if sdl::wcscmp(&text[0], copy) != 0 {
    return 3
  }
  return 0
}
`)
}

func TestSDL3OpaqueTypedefHandlesUseHeaderIdentity(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  context sdl::GLContext := sdl::gl_get_current_context()
  if context != nil {
    return 1
  }
  converter := sdl::iconv_open("UTF-8", "UTF-8")
  if sdl::iconv_close(converter) != 0 {
    return 2
  }
  return 0
}
`)
}

func TestSDL3HeaderFlagsAndHintConstantsHaveNativeValues(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  if sdl::WINDOW_RESIZABLE != 32 {
    return 1
  }
  if sdl::set_hint_with_priority(sdl::HINT_RENDER_DRIVER, "software", sdl::HintPriority::OVERRIDE) == false {
    return 2
  }
  if sdl::AUDIO_S16LE != 32784 {
    return 3
  }
  return 0
}
`)
}

func TestSDL3GamepadBindingNestedUnionsRetainTheirLayouts(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  binding := sdl::GamepadBinding{input_type: sdl::GamepadBindingType::BUTTON, input: sdl::GamepadBindingInput{button: 3}, output_type: sdl::GamepadBindingType::AXIS, output: sdl::GamepadBindingOutput{axis: sdl::GamepadBindingOutputAxis{axis: sdl::GamepadAxis::LEFTX, axis_min: -32768, axis_max: 32767}}}
  if binding.input.button != 3 {
    return 1
  }
  if binding.output.axis.axis != sdl::GamepadAxis::LEFTX {
    return 2
  }
  if binding.output.axis.axis_min != -32768 or binding.output.axis.axis_max != 32767 {
    return 3
  }
  return 0
}
`)
}

func TestSDL3X11HookAcceptsItsNullableCallback(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  callback sdl::X11EventHook := nil
  sdl::set_x11_event_hook(callback, nil)
  return 0
}
`)
}

func TestSDL3VulkanEntryPointIsNullableBeforeLoading(t *testing.T) {
	runSDLSource(t, `package main
use "vendor::sdl"
fn main() i32 {
  entry := sdl::vulkan_get_vk_get_instance_proc_addr()
  if entry != nil {
    return 1
  }
  return 0
}
`)
}

func TestSDL3ThreadMacroRunsCallbackAndReturnsItsStatus(t *testing.T) {
	runSDLSource(t, `package main
use "std::ffi::c"
use "vendor::sdl"
@[cc_macro="SDLCALL"]
fn worker(data rawptr) c::c_int {
  return 42
}
fn main() i32 {
  thread := sdl::create_thread(worker, "telia-worker", nil)
  if thread == nil {
    return 1
  }
  status c::c_int := 0
  sdl::wait_thread(thread, &status)
  if status != 42 {
    return 2
  }
  return 0
}
`)
}
