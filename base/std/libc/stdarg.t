package libc

#[default_cc="c"]
extern {
	#[link_name="llvm.va_start"]
  fn _va_start(arglist rawptr)

	#[link_name="llvm.va_end"]
  fn _va_end(arglist rawptr)

	#[link_name="llvm.va_copy"]
  fn _va_copy(dst rawptr, src rawptr)
}
