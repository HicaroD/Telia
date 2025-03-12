package io

#[default_cc="c"]
extern libc {
  fn printf(format string, args ...i32) i32
  fn puts(format string) i32
}

fn println(f string) {
  libc::puts(f)
  return
}
