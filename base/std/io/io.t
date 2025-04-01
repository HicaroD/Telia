package io

use "std::libc"

fn println(f string) {
  libc::println(f)
  return
}
