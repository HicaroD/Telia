package io

#[default_cc="c"]
extern libc {
  fn printf(format cstring, args ...f64) f64
  fn puts(format cstring) f64
}

fn println(f cstring) {
  libc::puts(f)
  return
}
