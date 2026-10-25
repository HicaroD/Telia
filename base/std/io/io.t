package io

@[default_cc="c"]
extern {
  @[link_name="_telia_println"]
  fn println(text string)
}
