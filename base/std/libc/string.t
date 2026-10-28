package libc

@[default_cc="c"]
extern strings {
  fn strcmp(str1 const cstring, str2 const cstring) c_int
  fn strlen(str const cstring) c_size_t
  fn strcat(str1 cstring, str2 const cstring) cstring
}
