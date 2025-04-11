package strings

use "std::libc"

#[default_cc="c"]
extern {
  #[link_name="strcmp"]
  fn compare(str1 @const string, str2 @const string) i32
}
