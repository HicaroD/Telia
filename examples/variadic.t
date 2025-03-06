package main

extern libc {
  fn printf(f cstring, @c args ...int)
}

fn main() {
  libc::printf("Hello, world!\n")
  return
}
