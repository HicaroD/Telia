package main

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
}

fn factorial(n u64) u64 {
  if n == 1 {
    return n
  }
  return n * factorial(n - 1)
}

fn main() i32 {
  result := factorial(6)
  libc.printf("Result: %d", result)
  return 0
}
