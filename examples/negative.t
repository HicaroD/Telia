package main

use "std::io"

fn main() i32 {
  negative := -10
  io::printf("POSITIVE: %d\n", -negative)
  io::printf("NEGATIVE: %d\n", negative)

  other_negative := 10
  io::printf("NEGATIVE: %d\n", -other_negative)
  io::printf("POSITIVE: %d\n", other_negative)
  return 0
}
