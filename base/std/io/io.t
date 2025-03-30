package io

use "std::libc"

fn println(f string) {
  libc::puts(f)
  return
}
