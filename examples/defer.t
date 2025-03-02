package main

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

#[linkage="external"]
fn test() {
  return
}

fn main() i32 {
  defer libc.puts("Hello, world 1")
  defer libc.puts("Hello, world 2")
  defer libc.puts("Hello, world 3")

  defer for i := 1; i < 10; i = i + 1 {
    libc.puts("From the if-defer")
  }

  return 0
}
