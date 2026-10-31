package main

use "std::libc"
use "std::io"
use "std::strings"
use "std::memory"

fn run() error {
  text := strings::clone("Hello") @catch err {
    return err
  }
  defer strings::free(&text)
  part := strings::slice(strings::view(&text), 0, 5) @catch err {
    return err
  }
  io::println(part)

  arena := memory::arena_create() @catch err {
    return err
  }
  defer memory::arena_destroy(&arena)
  greeting := strings::concat_in(&arena, part, ", Telia!") @catch err {
    return err
  }
  io::println(greeting)

  ctext := strings::to_c(greeting) @catch err {
    return err
  }
  defer strings::free_c(&ctext)
  libc::println(strings::c_view(&ctext))

  invalid := strings::slice(greeting, 5, 1) @catch err {
    io::println(err.msg)
    return nil
  }
  return error("expected invalid range")
}

fn main() i32 {
  run() @fail
  return 0
}
