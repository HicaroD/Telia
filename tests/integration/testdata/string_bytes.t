package main

use "std::io"

fn identity(text string) string {
  return text
}

fn main() i32 {
  text := identity("A\0B — Olá")
  io::println(text)
  io::println("")
  return 0
}
