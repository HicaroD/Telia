package libc

use "std::ffi::c"

@[default_cc="c"]
extern {
  fn getenv(name c::ConstCString) c::CString
  fn atoi(text c::ConstCString) c::c_int

  fn malloc(size c::c_size_t) rawptr
  fn calloc(num c::c_size_t, size c::c_size_t) rawptr
  fn realloc(ptr rawptr, new_size c::c_size_t) rawptr
  fn free(ptr rawptr)
  fn exit(code c::c_int)
}
