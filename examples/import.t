package main;

use "std::core";
use "package::something";

#[default_cc="c"]
extern C {
  fn puts(s cstring);
  fn printf(format cstring, ...);
}

fn main() {
  return;
}
