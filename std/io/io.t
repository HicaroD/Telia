package io

#[default_cc="c"]
extern libc {
  fn printf(format cstring, args ...f32) i32
  fn puts(format cstring) i32
}

fn printf(f cstring, args ...f32) {
  libc::printf(f, args)
  return
}

fn println(f cstring) {
  libc::puts(f)
  return
}
