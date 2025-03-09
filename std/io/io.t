package io

#[default_cc="c"]
extern libc {
  fn printf(format cstring, args ...i32) i32
  fn puts(format cstring) i32
}

fn println(f cstring) {
  libc::puts(f)
  return
}
