package main

use libc "std::libc"

fn succeed() error {
  return nil
}

fn main() i32 {
  succeed() @fail
  libc::printf("ok\n")
  return 0
}
