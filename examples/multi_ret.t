package main;

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn main() i32 {
  a := 1;
  libc.printf("%d", a);
  return 0;
}
