package main

use "std::io"
use "std::strings"

fn main() i32 {
  first := strings::clone("Hello\0é") @fail
  defer strings::free(&first)
  copy := strings::clone(strings::view(&first)) @fail
  if strings::view(&first) != strings::view(&copy) {
    return 1
  }
  strings::free(&first)
  strings::free(&first)
  io::println(strings::view(&copy))
  if strings::view(&first) == "" {
    io::println("reset")
  }
  joined := strings::concat(strings::view(&copy), "!") @fail
  defer strings::free(&joined)
  strings::free(&copy)
  io::println(strings::view(&joined))
  empty := strings::clone("") @fail
  strings::free(&empty)
  strings::free(&empty)
  return 0
}
