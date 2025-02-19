package main;

use "std::core";
use "package::something";

#[default_cc="c"]
extern C {
  fn puts(s *u8);
  fn printf(format *u8, ...);
}

fn main() {
  return;
}
