package main

use "std::io"

fn gcd(a int, b int) int {
  while b != 0 {
    if a > b {
      a = a - b
    }
    else {
      b = b - a
    }
  }
  return a
}

fn main() {
  b { name: "Hicaro" }
  a, b := 48, 18
  result := gcd(a, b)
  io::printf("GCD of %d and %d: %d", a, b, result)
  return
}
