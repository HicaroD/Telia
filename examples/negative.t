package main

use "std::libc"

fn main() i32 {
  negative := -10
  libc::printf("POSITIVE: %d\n", -negative)
  libc::printf("NEGATIVE: %d\n", negative)

  other_negative := 10
  libc::printf("NEGATIVE: %d\n", -other_negative)
  libc::printf("POSITIVE: %d\n", other_negative)
  return 0
}
