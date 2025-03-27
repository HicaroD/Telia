package main

use "std::io"

fn gcd(a i32, b i32) i32 {
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
  a i32, b i32 := 48, 18
  result := gcd(a, b)
  io::printf("GCD of %d and %d: %d\n", a, b, result)
  return
}
