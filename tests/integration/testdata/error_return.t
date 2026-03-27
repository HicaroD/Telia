package main

use libc "std::libc"

fn divide(a i32, b i32) (i32, error) {
  if b == 0 {
    return 0, error("division by zero")
  }
  return a / b, nil
}

fn main() i32 {
  val, err := divide(10, 2)
  libc::printf("%d\n", val)
  if err != nil {
    libc::printf("unexpected error\n")
  } else {
    libc::printf("ok\n")
  }

  val2, err2 := divide(10, 0)
  libc::printf("%d\n", val2)
  if err2 != nil {
    libc::printf("%s\n", err2.msg)
  }

  return 0
}
