package main

use "std::io"

fn identity(text string) string {
  return text
}

fn main() i32 {
  if identity("A\0B") == "A\0B" {
    io::println("equal")
  }
  if "A\0B" != "A\0C" {
    io::println("different")
  }
  if "" == "" {
    io::println("empty")
  }
  if "abc" != "ab" {
    io::println("length")
  }
  return 0
}
