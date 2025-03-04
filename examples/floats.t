package main

use "std::io"

fn main() {
  a := 2.0 + 1.0
  io::libc::printf("%f\n", a)
}
