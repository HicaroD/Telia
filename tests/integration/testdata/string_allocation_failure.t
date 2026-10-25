package main

use "std::strings"
use "std::memory"
use "std::io"

@[default_cc="c"]
extern {
  fn _test_fail_after(count i64)
  fn _test_live_allocations() i64
}

fn exercise() {
  block := memory::alloc(16) @fail
  defer memory::free(block)
  original := strings::clone("preserved") @fail
  defer strings::free(&original)
  _test_fail_after(0)
  zero, zero_err := memory::alloc(0)
  if zero_err == nil {
    if zero == nil {
      io::println("zero allocation")
    }
  }
  memory::free(zero)
  failed_block, alloc_err := memory::alloc(16)
  if alloc_err != nil {
    io::println("heap failed")
  }
  memory::free(failed_block)
  failed, err := strings::clone("new allocation")
  if err != nil {
    io::println("clone failed")
  }
  strings::free(&failed)
  joined, err2 := strings::concat(strings::view(&original), "suffix")
  if err2 != nil {
    io::println("concat failed")
  }
  strings::free(&joined)
  ctext, err3 := strings::to_c("")
  if err3 != nil {
    io::println("conversion failed")
  }
  strings::free_c(&ctext)
  io::println(strings::view(&original))
  _test_fail_after(-1)
  return
}

fn main() i32 {
  exercise()
  if _test_live_allocations() != 0 {
    return 1
  }
  io::println("released")
  return 0
}
