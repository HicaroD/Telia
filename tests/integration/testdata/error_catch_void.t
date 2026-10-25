package main
use "std::io"

use libc "std::libc"

fn failOnly() error {
  return error("boom")
}

fn main() i32 {
  failOnly() @catch err {
    io::println(err.msg)
    return 1
  }
  libc::printf("ok\n")
  return 0
}
