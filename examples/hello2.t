package main;

#[default_cc="c"]
extern C {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn main() i32 {
  message := "Hello, world";
  puts(message);
  return 0;
}
