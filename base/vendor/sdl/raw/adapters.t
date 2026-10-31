package raw

use "std::ffi::c"

@[cc="c"]
type NativeThreadEntry = fn(data rawptr) i32

struct NativeSize {
  width i32
  height i32
}

@[default_cc="c"]
extern {
  @[link_name="telia_sdl_render_debug_i32"]
  fn render_debug_i32(renderer *Renderer, x f32, y f32, value i32) bool

  @[link_name="telia_sdl_create_thread"]
  fn create_native_thread(entry NativeThreadEntry, name c::ConstCString, data rawptr) *Thread

  @[link_name="telia_sdl_io_from_bytes"]
  fn io_from_bytes_pointer(bytes *u8, size u64) *IOStream

  @[link_name="telia_sdl_read_bytes"]
  fn read_bytes_pointer(stream *IOStream, bytes *u8, size u64) u64
}

fn wait_native_thread(thread *Thread) i32 {
  status c::c_int := 0
  wait_thread(thread, &status)
  return status
}

fn window_size(window *Window, result *NativeSize) bool {
  width c::c_int := 0
  height c::c_int := 0
  if get_window_size(window, &width, &height) == false {
    return false
  }
  result.width = width
  result.height = height
  return true
}
