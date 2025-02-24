package main;

#[default_cc="c"]
extern C {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn print(message cstring) {
  C.puts(message);
  return;
}

fn main() i32 {
  message := "Hello, world";
  print(message);
  return 0;
}
