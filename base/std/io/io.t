package io

use "std::c"

fn println(f string) {
  c::puts(f)
  return
}
