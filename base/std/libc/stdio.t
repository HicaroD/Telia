package libc

#[default_cc="c"]
extern {
  fn printf(format @const string, args @c ...int) int

  #[link_name="puts"]
  fn println(format @const string) c_int
}
