package libc

@[default_cc="c"]
extern strings {
  fn strcmp(str1 @const cstring, str2 @const cstring) i32
  fn strlen(str @const cstring) u64
  fn strcat(str1 cstring, str2 @const cstring) cstring
}
