package libc

@[default_cc="c"]
extern {
  fn printf(format @const cstring, args @c ...int) int

  @[link_name="puts"]
  fn println(format @const cstring) c_int
}
