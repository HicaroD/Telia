package main

use libc "std::libc"

fn print_error(err error) {
  libc::printf("%s\n", err.msg)
}

fn make_error() error {
  return error("param broke")
}

fn main() i32 {
  print_error(make_error())
  return 0
}
