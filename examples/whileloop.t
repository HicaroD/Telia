package main

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

fn main() i32 {
  i := 0
  while i < 5 {
    libc.puts("Hello, world!")
    i = i + 1
  }
  return 0
}
