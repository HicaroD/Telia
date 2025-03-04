package main

use "std::io"
use "std::math"

fn main() {
  x := math::common::cos(0.0)
  if x == 1.0 {
    io::println("cos(0.0) = 1.0")
  }
  return
}
