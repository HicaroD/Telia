package main

#[default_cc="c"]
extern libc {

  fn printf(format cstring, ...) i32

  fn puts(format cstring) i32

}

#[linkage="external"]
fn test() {
  return
}

fn main() i32 {

  libc.puts("Hello, world!") @fail

  libc.puts("Hello, world!") @prop

  libc.puts("Hello, world!") @catch err {
    return 1
  }

  return 0
}
