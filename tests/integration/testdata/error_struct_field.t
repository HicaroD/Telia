package main

use libc "std::libc"

struct Result {
  err error
}

fn main() i32 {
  result := Result.{err: error("field broke")}
  libc::printf("%s\n", result.err.msg)
  return 0
}
