pkg main;

type Cstring = *u8;

#[default_cc="c"]
extern libc {
  fn printf(format Cstring, ...) i32;
  fn puts(format Cstring) i32;
}

fn print(message Cstring) {
  libc.puts(message);
  return;
}

fn main() i32 {
  libc.puts("Hello world");
  return 0;
}
