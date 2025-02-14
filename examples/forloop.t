pkg main;

extern libc {
  #[link_name="puts"] fn puts(format *u8) i32;
}

fn main() i32 {
  for i := 0; i <= 10; i = i + 1 {
    libc.puts("Hello, world");
  }
  return 0;
}
