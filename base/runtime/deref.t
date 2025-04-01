package runtime

use "std::libc"

fn _check_nil_pointer_deref(ptr rawptr) {
  if ptr == nil {
    libc::println("runtime panic: null pointer deref")
    libc::exit(1)
  }
  return
}
