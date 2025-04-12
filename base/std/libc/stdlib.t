package libc

@[default_cc="c"]
extern {
  fn malloc(size c_size_t) rawptr
  fn calloc(num c_size_t, size c_size_t) rawptr
  fn realloc(ptr rawptr, new_size c_size_t) rawptr
  fn free(ptr rawptr)
  fn exit(code c_int)
}
