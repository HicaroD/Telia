package sdl

use "vendor::sdl::raw"

type Size = raw::NativeSize

fn get_window_size(window *Window) (Size, error) {
  result := Size{width: 0, height: 0}
  if raw::window_size(window, &result) == false {
    return result, error("SDL get_window_size failed")
  }
  return result, nil
}
