package strings

use "std::libc"

#[default_cc="c"]
extern {
  #[link_name="strcmp"]
  fn compare(str1 string, str2 string) i32
}
