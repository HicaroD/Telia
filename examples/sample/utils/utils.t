package utils

use "std::io"

struct User {
  name string
}

fn foo() {
  io::println("from utils package")
  from_other_file_in_utils()
  return
}
