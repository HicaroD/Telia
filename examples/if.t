package main

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

fn main() i32 {
  if true {
    libc.puts("Hey, it is true!")
  } else {
    libc.puts("Hey, it is NOT true!")
  }
  libc.puts("Printing after conditions")
  return 0
}
