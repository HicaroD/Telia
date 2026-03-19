package main

use libc "std::libc"

struct Point {
  x int
  y int
}

fn sum_fields(p Point) int {
  return p.x + p.y
}

fn scale_x(p *Point, factor int) int {
  return p.x * factor
}

fn main() i32 {
  pt := Point.{
    x: 3,
    y: 7,
  }
  libc::printf("%d\n", pt.x)
  libc::printf("%d\n", pt.y)
  libc::printf("%d\n", sum_fields(pt))
  ptr := &pt
  libc::printf("%d\n", scale_x(ptr, 2))
  return 0
}
