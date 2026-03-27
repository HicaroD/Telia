package main

fn failOnly() error {
  return error("boom")
}

fn main() i32 {
  failOnly() @fail
  return 0
}
