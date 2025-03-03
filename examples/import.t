package main

#[default_cc="c"]
extern C {
  fn puts(s cstring)
  fn printf(format cstring, ...)
}

fn main() {
  return
}
