package main

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

fn print(message cstring) {
  libc.puts(message)
  return
}

fn main() i32 {
  print("Hello, world 😃")
  return 0
}
