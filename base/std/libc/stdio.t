package libc

#[default_cc="c"]
extern {
  fn printf(format @const string, args @c ...c_int) c_int
  fn puts(format @const string) c_int
}
