package main

use "std::libc"
use "std::io"

fn set(value *i32) {
  *value = 9
  return
}

fn compute() i32 {
  value i32 := 7
  defer set(&value)
  return value
}

fn nested() i32 {
  defer io::println("outer")
  if true {
    defer io::println("inner")
    return 1
  }
  return 0
}

fn main() i32 {
  libc::printf("%d\n", compute())
  value := nested()
  libc::printf("%d\n", value)
  return 0
}
