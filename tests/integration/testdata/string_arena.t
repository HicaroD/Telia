package main

use "std::strings"
use "std::memory"
use "std::io"
use "std::libc"

fn main() i32 {
  arena := memory::arena_create() @fail
  defer memory::arena_destroy(&arena)
  original := strings::clone_in(&arena, "stable\0é") @fail
  storage := memory::arena_alloc(&arena, 8192, 64) @fail
  io::println(original)
  joined := strings::concat_in(&arena, original, "!") @fail
  io::println(joined)
  ctext := strings::to_c_in(&arena, "C text") @fail
  libc::printf("%s\n", ctext)
  memory::arena_reset(&arena)
  reused := strings::clone_in(&arena, "new text") @fail
  io::println(reused)
  memory::arena_destroy(&arena)
  memory::arena_destroy(&arena)
  again := strings::clone_in(&arena, "again") @fail
  io::println(again)
  return 0
}
