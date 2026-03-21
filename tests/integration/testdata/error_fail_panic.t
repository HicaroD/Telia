package main

fn mightFail() (i32, error) {
  return 99, error("it failed")
}

fn main() i32 {
  val := mightFail() @fail
  return 0
}
