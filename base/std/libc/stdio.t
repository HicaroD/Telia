package libc

@[default_cc="c"]
extern {
  fn printf(format const cstring, args @c ...) c_int
  fn vprintf(format const cstring, args c_va_list) c_int

  @[link_name="puts"]
  fn println(format const cstring) c_int
}
