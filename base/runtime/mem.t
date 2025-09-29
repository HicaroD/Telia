package runtime

@[default_cc="c"]
extern {
  @[link_name="llvm.memcpy.p0.p0.i64"]
  fn memcpy(dest rawptr, src rawptr, size i64, isVolatile bool)
}
