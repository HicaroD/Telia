package main;

#[default_cc="c"]
extern C {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn no_return() {
  C.puts("No return");
  return;
}

fn main() i32 {
  no_return();
  return 0;
}
