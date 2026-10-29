package main

use "std::io"
use "std::strings"

struct Message {
  text string
}

fn pair(text string) (string, string) {
  return text, ""
}

fn main() i32 {
  zero := Message{}
  io::println(zero.text)
  if zero.text != "" {
    return 1
  }
  if strings::len(zero.text) != 0 {
    return 2
  }
  part := strings::slice(zero.text, 0, 0) @fail
  if part != "" {
    return 3
  }
  first, second := pair("A\0é")
  message := Message{text: first}
  copied := message
  io::println(copied.text)
  io::println(second)
  return 0
}
