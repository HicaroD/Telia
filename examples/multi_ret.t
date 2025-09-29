package main

use "std::io"

fn get(value i32) (i32, i32) {
  return 1, 2 + value
}

fn main() i32 {
  a, b := get(3)
  if a == 1 and b == 5 {
    io::println("a is 1 and b is 5")
  }

  c i32, d i64 := 10, 20
  e := 20
  if c == 10 and d == e {
    io::println("c is 10 and d is 20")
  }
  return 0
}
