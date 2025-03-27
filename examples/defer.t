package main

use "std::io"

fn main() i32 {
  defer io::println("Hello, world 1")
  defer io::println("Hello, world 2")
  defer io::println("Hello, world 3")

  defer for i i32 := 1; i < 10; i = i + 1 {
    io::println("From the if-defer")
  }
  return 0
}
