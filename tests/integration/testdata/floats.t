package main

use libc "std::libc"

fn add_f32(a f32, b f32) f32 {
  return a + b
}

fn scale_f64(x f64, factor f64) f64 {
  return x * factor
}

fn main() i32 {
  a f32 := 1.5 + 0.5
  libc::printf("%.1f\n", a)

  b f64 := 3.0 - 1.5
  libc::printf("%.1f\n", b)

  c f32 := add_f32(1.0, 2.0)
  libc::printf("%.1f\n", c)

  d f64 := scale_f64(2.5, 2.0)
  libc::printf("%.1f\n", d)

  e f64 := 9.0 / 4.0
  libc::printf("%.2f\n", e)

  if 1.0 < 2.0 {
    libc::printf("lt\n")
  }
  if 2.0 <= 2.0 {
    libc::printf("le\n")
  }
  if 3.0 > 2.0 {
    libc::printf("gt\n")
  }
  if 2.0 >= 2.0 {
    libc::printf("ge\n")
  }
  if 1.0 == 1.0 {
    libc::printf("eq\n")
  }
  if 1.0 != 2.0 {
    libc::printf("ne\n")
  }

  return 0
}
