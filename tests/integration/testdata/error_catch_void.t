package main

use libc "std::libc"

fn failOnly() error {
  return error("boom")
}

fn main() i32 {
  failOnly() @catch err {
    libc::printf("%s\n", err.msg)
    return 1
  }
  libc::printf("ok\n")
  return 0
}
