package main;

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn main() i32 {
  libc.puts("Hello, world!");
  libc.puts("Hello, world!");
  return 0;
}
