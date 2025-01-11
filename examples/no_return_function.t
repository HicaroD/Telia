extern C {
  fn printf(format *u8, ...) i32;
  fn puts(format *u8) i32;
}

fn no_return() {
  C.puts("No return");
  return;
}

fn main() i32 {
  no_return();
  return 0;
}
