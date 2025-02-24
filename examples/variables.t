package main;

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn main() i32 {
  name := "Hicaro";
  age, can_vote := 18, true;
  if can_vote {
    age, can_vote := 12, 10;
    libc.printf("Yes, you can vote because your age is %d", age);
  }
  return 0;
}
