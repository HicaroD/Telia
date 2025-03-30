package io

#[default_cc="c"]
extern {
  fn printf(format string, @c args ...i32) i32
  fn puts(format string) i32
}

fn println(f string) {
  puts(f)
  return
}
