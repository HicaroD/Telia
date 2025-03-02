package io

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

fn println(message cstring) {
  libc::print(message)
  libc::puts()
  return
}
