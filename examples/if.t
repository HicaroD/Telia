package main

use "std::io"

fn main() i32 {
  if true {
    io::libc::puts("Hey, it is true!")
  } 
  else {
    io::libc::puts("Hey, it is NOT true!")
  }
  io::libc::puts("Printing after conditions")
  return 0
}
