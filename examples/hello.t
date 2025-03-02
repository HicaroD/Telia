package main

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

fn main() i32 {
  libc.printf("Hello, world!\n") @fail
  libc.printf("Hello, world!\n") @prop
  libc.printf("Hello, world!\n") @catch err {
    return 1
  }
  return 0
}
