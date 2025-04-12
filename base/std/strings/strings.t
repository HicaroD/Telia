package strings

#[default_cc="c"]
extern {
  #[link_name="strcmp"]
  fn compare(str1 @const string, str2 @const string) i32

  #[link_name="strlen"]
  fn len(str @const string) i64
}
