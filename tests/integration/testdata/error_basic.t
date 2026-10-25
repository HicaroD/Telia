package main
use "std::io"

use libc "std::libc"

fn main() i32 {
  e := error("oops")
  io::println(e.msg)
  return 0
}
