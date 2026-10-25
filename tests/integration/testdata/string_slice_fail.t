package main

use "std::strings"

fn main() i32 {
  strings::slice("a", 2, 1) @fail
  return 0
}
