extern C {
  fn printf(format *u8, ...) i32;
  fn puts(format *u8) i32;
}

fn main() i32 {
  if true {
    puts("Hey, it is true!");
  } else {
    puts("Hey, it is NOT true!");
  }
  puts("Printing after conditions");
  return 0;
}
