package main

use "std::libc"
use "std::strings"
use "vendor::sdl"

@[default_cc="c"]
extern {
  fn getenv(name const cstring) cstring
  fn atoi(text const cstring) i32
}

fn main() i32 {
  frame_limit i32 := 0
  configured_limit := getenv("TELIA_SDL_EXAMPLE_FRAMES")
  if configured_limit != nil {
    frame_limit = atoi(configured_limit)
    if frame_limit <= 0 {
      libc::println("TELIA_SDL_EXAMPLE_FRAMES must be a positive integer")
      return 9
    }
  }
  if sdl::init() == false {
    libc::println(sdl::get_error())
    return 1
  }
  defer sdl::quit()

  title := strings::to_c("Telia + SDL3") @fail
  defer strings::free_c(&title)
  window := sdl::create_window(strings::c_view(&title), 640, 480)
  if window == nil {
    libc::println(sdl::get_error())
    return 2
  }
  defer sdl::destroy_window(window)
  renderer := sdl::create_renderer(window)
  if renderer == nil {
    libc::println(sdl::get_error())
    return 3
  }
  defer sdl::destroy_renderer(renderer)

  if sdl::raise_window(window) == false {
    libc::println(sdl::get_error())
    return 10
  }

  frame i32 := 0
  while frame_limit == 0 or frame < frame_limit {
    event := sdl::poll_event()
    while event != 0 {
      if sdl::is_quit_event(event) {
        return 0
      }
      event = sdl::poll_event()
    }
    if sdl::set_draw_color(renderer, 20, 24, 36, 255) == false {
      libc::println(sdl::get_error())
      return 4
    }
    if sdl::clear(renderer) == false {
      libc::println(sdl::get_error())
      return 5
    }
    if sdl::set_draw_color(renderer, 90, 180, 255, 255) == false {
      libc::println(sdl::get_error())
      return 6
    }
    if sdl::fill_rect(renderer, 220.0, 160.0, 200.0, 160.0) == false {
      libc::println(sdl::get_error())
      return 7
    }
    if sdl::present(renderer) == false {
      libc::println(sdl::get_error())
      return 8
    }
    sdl::delay(16)
    frame = frame + 1
  }
  return 0
}
