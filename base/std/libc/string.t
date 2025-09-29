package libc

@[default_cc="c"]
extern strings {
  fn strcmp(str1 @const string, str2 @const string) i32
  fn strlen(str @const string) i64
  fn strcat(str1 @const string, str2 @const string) string
}
