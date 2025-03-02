package main;

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn get() (int, int) {
    return 1, 1;
}

fn main() i32 {
  a, b := get();
  return 0;
}