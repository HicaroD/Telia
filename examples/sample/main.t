package main;

#[default_cc="c"]
extern libc {
  fn printf(format *u8, ...) i32;
  fn puts(format *u8) i32;
}

fn print(message *u8) {
  libc.puts(message);
  return;
}

fn main() i32 {
  print("Hello, world 😃");
  return 0;
}
