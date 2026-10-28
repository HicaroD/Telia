package sdl

@[default_cc="c"]
extern {
  @[link_name="telia_sdl_init"]
  fn init() bool
  @[link_name="telia_sdl_quit"]
  fn quit()
  @[link_name="telia_sdl_get_error"]
  fn get_error() cstring
  @[link_name="telia_sdl_create_window"]
  fn create_window(title cstring, width i32, height i32) rawptr
  @[link_name="telia_sdl_destroy_window"]
  fn destroy_window(window rawptr)
  @[link_name="telia_sdl_raise_window"]
  fn raise_window(window rawptr) bool
  @[link_name="telia_sdl_create_renderer"]
  fn create_renderer(window rawptr) rawptr
  @[link_name="telia_sdl_destroy_renderer"]
  fn destroy_renderer(renderer rawptr)
  @[link_name="telia_sdl_poll_event"]
  fn poll_event() u32
  @[link_name="telia_sdl_is_quit_event"]
  fn is_quit_event(event u32) bool
  @[link_name="telia_sdl_set_draw_color"]
  fn set_draw_color(renderer rawptr, r u8, g u8, b u8, a u8) bool
  @[link_name="telia_sdl_clear"]
  fn clear(renderer rawptr) bool
  @[link_name="telia_sdl_fill_rect"]
  fn fill_rect(renderer rawptr, x f32, y f32, width f32, height f32) bool
  @[link_name="telia_sdl_present"]
  fn present(renderer rawptr) bool
  @[link_name="telia_sdl_delay"]
  fn delay(milliseconds u32)
}
