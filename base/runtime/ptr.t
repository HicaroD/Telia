package runtime

use "std::c"
use "std::io"

fn __check_nil_pointer_deref(ptr rawptr) {
  c.puts("check nil deref from runtime package")
  exit(1)
  return
}
