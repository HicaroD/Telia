package libc

use "std::ffi::c"

@[default_cc="c"]
extern strings {
  fn strcmp(str1 c::ConstCString, str2 c::ConstCString) c::c_int
  fn strlen(str c::ConstCString) c::c_size_t
  fn strcat(str1 c::CString, str2 c::ConstCString) c::CString
}
