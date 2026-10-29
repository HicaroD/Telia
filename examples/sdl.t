package main

use "std::libc"
use "std::strings"
use "vendor::sdl"

fn main() i32 {
  frame_limit i32 := 0
  configured_limit := libc::getenv("TELIA_SDL_EXAMPLE_FRAMES")
  if configured_limit != nil {
    frame_limit = libc::atoi(configured_limit)
    if frame_limit <= 0 {
      libc::println("TELIA_SDL_EXAMPLE_FRAMES must be a positive integer")
      return 9
    }
  }
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    libc::println(sdl::get_error())
    return 1
  }
  defer sdl::quit()

  title := strings::to_c("Telia + SDL3") @fail
  defer strings::free_c(&title)
  window := sdl::create_window(strings::c_view(&title), 640, 480, 0)
  if window == nil {
    libc::println(sdl::get_error())
    return 2
  }
  defer sdl::destroy_window(window)
  renderer := sdl::create_renderer(window, nil)
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
  event := sdl::Event.{type: 0}
  while frame_limit == 0 or frame < frame_limit {
    while sdl::poll_event(&event) {
      if event.type == sdl::EVENT_QUIT {
        return 0
      }
    }
    if sdl::set_render_draw_color(renderer, 20, 24, 36, 255) == false {
      libc::println(sdl::get_error())
      return 4
    }
    if sdl::render_clear(renderer) == false {
      libc::println(sdl::get_error())
      return 5
    }
    if sdl::set_render_draw_color(renderer, 90, 180, 255, 255) == false {
      libc::println(sdl::get_error())
      return 6
    }
    rect := sdl::FRect.{x: 220.0, y: 160.0, w: 200.0, h: 160.0}
    if sdl::render_fill_rect(renderer, &rect) == false {
      libc::println(sdl::get_error())
      return 7
    }
    if sdl::render_present(renderer) == false {
      libc::println(sdl::get_error())
      return 8
    }
    sdl::delay(16)
    frame = frame + 1
  }
  return 0
}
