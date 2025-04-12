package main

use "std::strings"
use "std::io"

fn main() {
  size := strings::len("HELLO")
  if size == 5 {
    io::println("size of hello is 5")
  }
  
  if strings::compare("Hello", "Hello") != 0 {
    io::println("strings are different")
  }
}
