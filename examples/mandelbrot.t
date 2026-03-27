package main

use "std::io"
use "std::libc"

fn mandelbrot(cr f64, ci f64) i32 {
  zr f64 := 0.0
  zi f64 := 0.0
  n i32 := 0
  two f64 := 2.0
  four f64 := 4.0
  escaped i32 := 0
  escapeIter i32 := 0
  while n < 80 {
    if escaped == 0 {
      tr f64 := zr * zr - zi * zi + cr
      ti f64 := two * zr * zi + ci
      zr = tr
      zi = ti
    }
    if zr * zr + zi * zi > four {
      if escaped == 0 {
        escaped = 1
        escapeIter = n
      }
    }
    n = n + 1
  }
  if escaped == 0 {
    return 0
  }
  return escapeIter
}

fn pixel(m i32) {
  if m == 0 {
    libc::printf("%c[38;5;15m@", 27)
  } elif m < 3 {
    libc::printf("%c[38;5;196m@", 27)
  } elif m < 6 {
    libc::printf("%c[38;5;208m@", 27)
  } elif m < 12 {
    libc::printf("%c[38;5;226m@", 27)
  } elif m < 25 {
    libc::printf("%c[38;5;46m@", 27)
  } elif m < 50 {
    libc::printf("%c[38;5;51m@", 27)
  } else {
    libc::printf("%c[38;5;21m@", 27)
  }
}

fn main() i32 {
  w f64 := 100.0
  h f64 := 40.0
  three f64 := 3.0
  half f64 := 2.0
  y f64 := 0.0
  while y < 40.0 {
    x f64 := 0.0
    while x < 100.0 {
      cr f64 := x / w * three - 2.0
      ci f64 := y / h * half - 1.0
      m i32 := mandelbrot(cr, ci)
      pixel(m)
      x = x + 1.0
    }
    libc::printf("%c[0m\n", 27)
    y = y + 1.0
  }
  return 0
}
