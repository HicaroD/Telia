package libc

#[default_cc="c"]
extern {
  fn exit(code c_int)
}
