package main

@[default_cc="c"]
extern libc {
  fn printf(format cstring, ...) i32
  fn puts(format cstring) i32
}

fn connect_to_db() (int, error) {
    return 1, error("db connection failed")
}

fn safe_connect() error {
    return nil
}

fn main() i32 {
  // @fail on error-returning function: panics if error is non-nil
  safe_connect() @fail

  // @fail on tuple-returning function: extracts non-error values
  db := connect_to_db() @fail
  libc::printf("%d\n", db)
  return 0
}
