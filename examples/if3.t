extern libc {
  fn puts(format *u8) i32;
}

fn main() i32 {
  can_vote := true;
  if can_vote {
    other_bool := true;
    if other_bool {
      libc.puts("Hey, it is true and other_bool as well!");
    }
    libc.puts("Hey, it is true!");
  } else {
    libc.puts("Hey, it is NOT true!");
  }
  return 0;
}
