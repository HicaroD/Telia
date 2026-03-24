package main

use libc "std::libc"

fn connect(fail bool) (i32, error) {
  if fail {
    return 0, error("connection refused")
  }
  return 42, nil
}

fn main() i32 {
  val := connect(false) @catch err {
    libc::printf("error: %s\n", err.msg)
    return 1
  }
  libc::printf("val: %d\n", val)
  return 0
}
