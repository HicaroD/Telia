package c

#[default_cc="c"]
extern _ {
	#[link_name="llvm.va_start"] fn _va_start(arglist *i8)
	#[link_name="llvm.va_end"]   fn _va_end(arglist *i8)
	#[link_name="llvm.va_copy"]  fn _va_copy(dst *i8, src *i8)
}
