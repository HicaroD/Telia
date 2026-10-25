package main

use "std::io"
use "std::libc"
use "std::strings"

fn main() i32 {
  converted := strings::to_c("Olá") @fail
  defer strings::free_c(&converted)
  pointer := strings::c_view(&converted)
  libc::printf("%s\n", pointer)
  io::println(strings::from_c(pointer))
  empty := strings::to_c("") @fail
  defer strings::free_c(&empty)
  libc::printf("[%s]\n", strings::c_view(&empty))
  rejected, err := strings::to_c("a\0b")
  if err != nil {
    io::println("rejected")
  }
  strings::free_c(&rejected)
  io::println(strings::from_c("static"))
  return 0
}
