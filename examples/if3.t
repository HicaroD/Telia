package main

use "std::io"

fn main() i32 {
  can_vote := true
  if can_vote {
    other_bool := true
    if other_bool {
      io::puts("Hey, can_vote is true and other_bool as well!")
    }
    io::puts("Hey, can_vote is true!")
  } else {
    io::puts("Hey, can_vote is NOT true!")
  }
  return 0
}
