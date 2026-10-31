package main

use "std::libc"
use "std::io"
use "vendor::sdl"
use "pkg::game"

fn rectangle(renderer *sdl::Renderer, x f32, y f32, w f32, h f32, r u8, g u8, b u8) bool {
  if sdl::set_render_draw_color(renderer, r, g, b, 255) == false {
    return false
  }
  rect := sdl::FRect{x: x, y: y, w: w, h: h}
  return sdl::render_fill_rect(renderer, &rect)
}

fn label(renderer *sdl::Renderer, x f32, y f32, text string) bool {
  if sdl::set_render_scale(renderer, 2.0, 2.0) == false {
    return false
  }
  drawn := sdl::render_debug_text(renderer, x / 2.0, y / 2.0, text)
  restored := sdl::set_render_scale(renderer, 1.0, 1.0)
  return drawn == nil and restored
}

fn cell_position(cell i32, origin f32) f32 {
  index i32 := 0
  result := origin
  while index < cell {
    result = result + 20.0
    index = index + 1
  }
  return result
}

fn draw(renderer *sdl::Renderer, state *game::State, best i32) bool {
  if sdl::set_render_draw_color(renderer, 12, 17, 25, 255) == false {
    return false
  }
  if sdl::render_clear(renderer) == false {
    return false
  }
  if rectangle(renderer, 24.0, 88.0, 560.0, 480.0, 18, 27, 37) == false {
    return false
  }
  if rectangle(renderer, 608.0, 88.0, 288.0, 480.0, 18, 27, 37) == false {
    return false
  }
  if sdl::set_render_draw_color(renderer, 25, 36, 47, 255) == false {
    return false
  }
  x f32 := 24.0
  while x <= 584.0 {
    if sdl::render_line(renderer, x, 88.0, x, 568.0) == false {
      return false
    }
    x = x + 20.0
  }
  y f32 := 88.0
  while y <= 568.0 {
    if sdl::render_line(renderer, 24.0, y, 584.0, y) == false {
      return false
    }
    y = y + 20.0
  }
  if state.won == false {
    food_x := cell_position(state.food.x, 24.0)
    food_y := cell_position(state.food.y, 88.0)
    if rectangle(renderer, food_x + 3.0, food_y + 3.0, 14.0, 14.0, 255, 153, 92) == false {
      return false
    }
    if rectangle(renderer, food_x + 7.0, food_y + 1.0, 4.0, 4.0, 150, 222, 123) == false {
      return false
    }
  }
  index := state.length
  while index > 0 {
    index = index - 1
    cell := state.body[index]
    body_x := cell_position(cell.x, 24.0)
    body_y := cell_position(cell.y, 88.0)
    if rectangle(renderer, body_x + 1.0, body_y + 1.0, 18.0, 18.0, 55, 173, 136) == false {
      return false
    }
  }
  head := state.body[0]
  head_x := cell_position(head.x, 24.0)
  head_y := cell_position(head.y, 88.0)
  if rectangle(renderer, head_x + 1.0, head_y + 1.0, 18.0, 18.0, 119, 235, 179) == false {
    return false
  }
  eye_x f32 := 5.0
  eye_y f32 := 5.0
  other_x f32 := 12.0
  other_y f32 := 5.0
  if state.direction == game::Direction::Down {
    eye_y = 12.0
    other_y = 12.0
  } elif state.direction == game::Direction::Left {
    other_x = 5.0
    other_y = 12.0
  } elif state.direction == game::Direction::Right {
    eye_x = 12.0
    other_y = 12.0
  }
  if rectangle(renderer, head_x + eye_x, head_y + eye_y, 3.0, 3.0, 12, 35, 35) == false {
    return false
  }
  if rectangle(renderer, head_x + other_x, head_y + other_y, 3.0, 3.0, 12, 35, 35) == false {
    return false
  }
  if sdl::set_render_draw_color(renderer, 119, 235, 179, 255) == false {
    return false
  }
  if label(renderer, 24.0, 24.0, "TELIA / SNAKE") == false {
    return false
  }
  if sdl::set_render_draw_color(renderer, 143, 161, 178, 255) == false {
    return false
  }
  if sdl::render_debug_text(renderer, 24.0, 60.0, "A LITTLE GAME. A NEW LANGUAGE. SDL3.") != nil {
    return false
  }
  if label(renderer, 624.0, 112.0, "SCORE") == false {
    return false
  }
  if sdl::set_render_draw_color(renderer, 229, 239, 245, 255) == false {
    return false
  }
  if sdl::set_render_scale(renderer, 2.0, 2.0) == false {
    return false
  }
  score_drawn := sdl::render_debug_i32(renderer, 312.0, 76.0, state.score)
  best_drawn := sdl::render_debug_i32(renderer, 312.0, 132.0, best)
  scale_restored := sdl::set_render_scale(renderer, 1.0, 1.0)
  if score_drawn != nil or best_drawn != nil or scale_restored == false {
    return false
  }
  if sdl::set_render_draw_color(renderer, 143, 161, 178, 255) == false {
    return false
  }
  if label(renderer, 624.0, 224.0, "SESSION BEST") == false {
    return false
  }
  if sdl::render_debug_text(renderer, 624.0, 344.0, "ARROWS / WASD\n\nMOVE\n\n\nSPACE\n\nPAUSE\n\n\nR  RESTART\n\nESC  QUIT") != nil {
    return false
  }
  if sdl::render_debug_text(renderer, 24.0, 592.0, "EAT. GROW. DON'T HIT THE WALLS OR YOURSELF.") != nil {
    return false
  }
  if state.paused or state.over or state.won {
    if rectangle(renderer, 152.0, 258.0, 304.0, 128.0, 12, 17, 25) == false {
      return false
    }
    if sdl::set_render_draw_color(renderer, 229, 239, 245, 255) == false {
      return false
    }
    if state.over {
      if label(renderer, 232.0, 282.0, "GAME OVER") == false {
        return false
      }
    } elif state.won {
      if label(renderer, 240.0, 282.0, "YOU WIN!") == false {
        return false
      }
    } else {
      if label(renderer, 256.0, 282.0, "PAUSED") == false {
        return false
      }
    }
    if state.paused {
      if sdl::render_debug_text(renderer, 232.0, 338.0, "SPACE TO CONTINUE") != nil {
        return false
      }
    } else {
      if sdl::render_debug_text(renderer, 256.0, 338.0, "R TO RESTART") != nil {
        return false
      }
    }
  }
  return sdl::render_present(renderer)
}

fn main() i32 {
  frame_limit i32 := 0
  configured := libc::getenv("TELIA_SNAKE_FRAMES")
  if configured != nil {
    frame_limit = libc::atoi(configured)
    if frame_limit <= 0 {
      io::println("TELIA_SNAKE_FRAMES must be a positive integer")
      return 1
    }
  }
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    io::println(sdl::get_error())
    return 2
  }
  defer sdl::quit()
  window := sdl::create_window("Telia / Snake", 920, 624, 0) @catch err {
    io::println(err.msg)
    return 3
  }
  defer sdl::destroy_window(window)
  renderer := sdl::create_renderer(window) @catch err {
    io::println(err.msg)
    return 4
  }
  defer sdl::destroy_renderer(renderer)
  if sdl::set_render_logical_presentation(renderer, 920, 624, sdl::RendererLogicalPresentation::LETTERBOX) == false {
    io::println(sdl::get_error())
    return 5
  }
  if sdl::raise_window(window) == false {
    io::println(sdl::get_error())
    return 6
  }
  state := game::create(sdl::rand_bits())
  best i32 := 0
  event := sdl::Event{type: 0}
  frame i32 := 0
  previous := sdl::get_ticks()
  accumulated u64 := 0
  while frame_limit == 0 or frame < frame_limit {
    now := sdl::get_ticks()
    elapsed := now - previous
    previous = now
    if elapsed > 250 {
      elapsed = 250
    }
    while sdl::poll_event(&event) {
      if event.type == sdl::EVENT_QUIT {
        return 0
      }
      if event.type == sdl::EVENT_KEY_DOWN and event.key.repeat == false {
        code := event.key.scancode
        if code == sdl::Scancode::ESCAPE {
          return 0
        } elif code == sdl::Scancode::UP or code == sdl::Scancode::W {
          game::turn(&state, game::Direction::Up)
        } elif code == sdl::Scancode::DOWN or code == sdl::Scancode::S {
          game::turn(&state, game::Direction::Down)
        } elif code == sdl::Scancode::LEFT or code == sdl::Scancode::A {
          game::turn(&state, game::Direction::Left)
        } elif code == sdl::Scancode::RIGHT or code == sdl::Scancode::D {
          game::turn(&state, game::Direction::Right)
        } elif code == sdl::Scancode::SPACE {
          game::toggle_pause(&state)
          accumulated = 0
          elapsed = 0
        } elif code == sdl::Scancode::R {
          state = game::create(sdl::rand_bits())
          accumulated = 0
          elapsed = 0
        }
      }
    }
    if state.paused == false and state.over == false and state.won == false {
      accumulated = accumulated + elapsed
      interval u64 := 120
      if state.score >= 20 {
        interval = 70
      } elif state.score >= 10 {
        interval = 90
      }
      while accumulated >= interval {
        game::step(&state)
        accumulated = accumulated - interval
      }
    } else {
      accumulated = 0
    }
    if state.score > best {
      best = state.score
    }
    if draw(renderer, &state, best) == false {
      io::println(sdl::get_error())
      return 7
    }
    sdl::delay(16)
    frame = frame + 1
  }
  return 0
}
