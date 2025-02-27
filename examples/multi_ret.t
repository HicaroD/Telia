package main;

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn get(value i32) (i32, i32) {
  return 1, 2 + value;
}

fn main() i32 {
  a, b, c, d, e := get(1), 1, get(2);
  libc.printf("%d %d %d %d %d", a, b, c, d, e);
  return 0;
}
