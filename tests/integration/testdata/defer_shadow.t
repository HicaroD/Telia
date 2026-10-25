package main
use "std::libc"

fn show(value *i32) {
  libc::printf("%d\n", *value)
  return
}

fn nested() i32 {
  value i32 := 1
  defer show(&value)
  if true {
    value i32 := 2
    defer show(&value)
    return 0
  }
  return 0
}

fn main() i32 {
  return nested()
}
