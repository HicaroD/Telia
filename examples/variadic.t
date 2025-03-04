package main

extern libc {
  fn printf(f cstring, @for_c args ...int)
}

fn main() {
  libc::printf("hello", 1, 2)
  return
}
