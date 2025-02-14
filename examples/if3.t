pkg main;

extern libc {
  fn puts(format *u8) i32;
}

fn main() i32 {
  can_vote := true;
  if can_vote {
    other_bool := true;
    if other_bool {
      libc.puts("Hey, can_vote is true and other_bool as well!");
    }
    libc.puts("Hey, can_vote is true!");
  } else {
    libc.puts("Hey, can_vote is NOT true!");
  }
  return 0;
}
