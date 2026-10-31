package sdl

use "vendor::sdl::raw"
use "std::strings"

struct OptionalString {
  value string
  present bool
}

fn create_renderer(window *Window) (*Renderer, error) {
  result := raw::create_renderer(window, nil)
  if result == nil {
    return nil, error("SDL create_renderer failed")
  }
  return result, nil
}

fn create_renderer_named(window *Window, name string) (*Renderer, error) {
  converted, err := strings::to_c(name)
  if err != nil {
    return nil, err
  }
  defer strings::free_c(&converted)
  result := raw::create_renderer(window, strings::c_view(&converted))
  if result == nil {
    return nil, error("SDL create_renderer_named failed")
  }
  return result, nil
}

fn render_debug_i32(renderer *Renderer, x f32, y f32, value i32) error {
  if raw::render_debug_i32(renderer, x, y, value) == false {
    return error("SDL render_debug_i32 failed")
  }
  return nil
}
