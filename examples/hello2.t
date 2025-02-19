package main;

#[default_cc="c"]
extern C {
  fn printf(format *u8, ...) i32;
  fn puts(format *u8) i32;
}

fn main() i32 {
  message := "Hello, world";
  puts(message);
  return 0;
}
