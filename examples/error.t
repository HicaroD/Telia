package main;

#[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32;
  fn puts(format cstring) i32;
}

fn connect_to_db() (int, error) {
    return 1, None;
}

fn main() i32 {
  db := connect_to_db() @catch err {
    return err;
  };

  db := connect_to_db() @prop;

  db := connect_to_db() @panic;
  return 0;
}