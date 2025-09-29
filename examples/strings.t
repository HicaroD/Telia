package main

use "std::libc"
use "std::io"

fn main() {
  size := libc::strings::strlen("HELLO")
  if size == 5 {
    io::println("size of hello is 5")
  }
  
  if libc::strings::strcmp("Hello", "Hello") == 0 {
    io::println("equal")
  }
}
