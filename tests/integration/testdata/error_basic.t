package main

use libc "std::libc"

fn main() i32 {
  e := error("oops")
  libc::printf("%s\n", e.msg)
  return 0
}
