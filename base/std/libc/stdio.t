package libc

use "std::ffi::c"

@[default_cc="c"]
extern {
  fn printf(format c::ConstCString, args @c ...) c::c_int
  fn vprintf(format c::ConstCString, args c::c_va_list) c::c_int
  fn puts(text c::ConstCString) c::c_int

  @[link_name="puts"]
  fn println(format c::ConstCString) c::c_int
}
