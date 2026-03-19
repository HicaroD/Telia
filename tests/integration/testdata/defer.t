package main

use libc "std::libc"

fn lifo() i32 {
  defer libc::printf("first\n")
  defer libc::printf("second\n")
  defer libc::printf("third\n")
  return 0
}

fn no_defer() i32 {
  return 42
}

fn nested(x i32) i32 {
  if x > 0 {
    defer libc::printf("inner\n")
  }
  defer libc::printf("outer\n")
  return 0
}

fn main() i32 {
  lifo()
  no_defer()
  nested(1)
  return 0
}
