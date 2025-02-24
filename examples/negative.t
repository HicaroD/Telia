package main;

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32;
}

fn main() i32 {
  negative := -10;
  libc.printf("%d", -negative);
  libc.printf("%d", negative);

  other_negative := 10;
  libc.printf("%d", -other_negative);
  libc.printf("%d", other_negative);
  return 0;
}
