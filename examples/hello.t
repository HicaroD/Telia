#[default_cc="c"]
extern libc {
  fn printf(format *u8, ...) i32;
  fn puts(format *u8) i32;
  fn islower(c int) int;
}

fn print(message *u8) {
  libc.puts(message);
  return;
}

fn main() i32 {
  libc.puts("Hello world");

  number := 98;
  is_lower := libc.islower(number);
  if is_lower != 0 {
    libc.printf("%d is lowercase", number);
  }

  return 0;
}
