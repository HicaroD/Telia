package main

use "std::io"

fn main() i32 {
  age, can_vote := 18, true
  if can_vote {
    my_age, can_vote := 12, 10
    io::printf("Yes, you can vote because your age is %d %d", my_age, age)
  }
  return 0
}
