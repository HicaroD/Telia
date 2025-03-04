package main

extern libc {
  fn printf(f cstring, @for_c args ...int)
}

fn main() {
  libc::printf("Hello, world!\n")
  return
}
