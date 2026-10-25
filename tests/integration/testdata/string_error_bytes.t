package main

fn fail() error {
  return error("A\0Bé")
}

fn main() i32 {
  fail() @fail
  return 0
}
