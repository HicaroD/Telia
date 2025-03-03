package io

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

fn printf(f cstring, ...) {
  libc::printf(f)
  return
}

fn println(f cstring) {
  libc::puts(f)
  return
}
