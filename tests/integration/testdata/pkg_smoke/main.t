package main

use "std::io"
use "pkg::greet"

fn main() i32 {
  greet::hello()
  io::println("Hello from main package!")
  return 0
}
