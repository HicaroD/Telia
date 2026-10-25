package main

use "std::io"
use "std::strings"

fn main() i32 {
  part := strings::slice("A\0Bé", 1, 3) @fail
  io::println(part)
  byte := strings::slice("é", 0, 1) @fail
  io::println(byte)
  empty := strings::slice("", 0, 0) @fail
  if strings::len(empty) == 0 {
    io::println("empty")
  }
  if strings::len("é") == 2 {
    io::println("bytes")
  }
  reversed, err := strings::slice("abc", 2, 1)
  if err != nil {
    io::println("reversed")
  }
  past_end, err2 := strings::slice("abc", 0, 4)
  if err2 != nil {
    io::println("past end")
  }
  return 0
}
