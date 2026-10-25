package main

use "std::memory"
use "std::strings"
use "std::io"

@[default_cc="c"]
extern {
  fn _test_fail_after(count i64)
  fn _test_live_allocations() i64
  fn _test_is_aligned(ptr rawptr, alignment u64) bool
}

fn main() i32 {
  arena := memory::arena_create() @fail
  if _test_live_allocations() != 0 {
    return 1
  }
  original := strings::clone_in(&arena, "preserved") @fail
  aligned := memory::arena_alloc(&arena, 7, 64) @fail
  if _test_is_aligned(aligned, 64) {
  } else {
    return 4
  }
  _test_fail_after(0)
  failed, err := memory::arena_alloc(&arena, 8192, 64)
  if err != nil {
    io::println("growth failed")
  }
  if failed != nil {
    return 2
  }
  io::println(original)
  memory::arena_reset(&arena)
  reused := strings::clone_in(&arena, "reused") @fail
  io::println(reused)
  overflow, err2 := memory::arena_alloc(&arena, 18446744073709551615, 64)
  if err2 != nil {
    io::println("overflow rejected")
  }
  invalid, err3 := memory::arena_alloc(&arena, 1, 3)
  if err3 != nil {
    io::println("alignment rejected")
  }
  memory::arena_destroy(&arena)
  memory::arena_destroy(&arena)
  if _test_live_allocations() != 0 {
    return 3
  }
  io::println("released")
  return 0
}
