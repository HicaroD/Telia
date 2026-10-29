package main
use "std::io"

use libc "std::libc"

struct Result {
  err error
}

fn main() i32 {
  result := Result{err: error("field broke")}
  err := result.err
  io::println(err.msg)
  return 0
}
