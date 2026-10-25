package main
use "std::io"

use libc "std::libc"

fn print_error(err error) {
  io::println(err.msg)
}

fn make_error() error {
  return error("param broke")
}

fn main() i32 {
  print_error(make_error())
  return 0
}
