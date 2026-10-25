package main

use "std::libc"
use "std::io"
use "std::strings"
use "std::memory"

fn main() i32 {
  text := strings::clone("Hello") @fail
  defer strings::free(&text)
  part := strings::slice(strings::view(&text), 0, 5) @fail
  io::println(part)

  arena := memory::arena_create() @fail
  defer memory::arena_destroy(&arena)
  greeting := strings::concat_in(&arena, part, ", Telia!") @fail
  io::println(greeting)

  ctext := strings::to_c(greeting) @fail
  defer strings::free_c(&ctext)
  libc::println(strings::c_view(&ctext))

  invalid := strings::slice(greeting, 5, 1) @catch err {
    io::println(err.msg)
    return 0
  }
  return 1
}
