package main
use "std::io"

fn fail() error {
  return error("")
}

fn main() i32 {
  err := fail()
  if err != nil {
    io::println("present")
  }
  fail() @catch err2 {
    io::println(err2.msg)
    io::println("caught")
  }
  return 0
}
