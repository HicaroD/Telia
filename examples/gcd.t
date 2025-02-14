pkg main;

#[default_cc="c"]
extern C {
  fn puts(s *u8);
  fn printf(format *u8, ...);
}

fn gcd(a int, b int) int {
  while b != 0 {
    if a > b {
      a = a - b;
    }
    else {
      b = b - a;
    }
  }
  return a;
}

fn main() {
  a, b := 48, 18;
  result := gcd(a, b);
  C.printf("GCD of %d and %d: %d", a, b, result);
  return;
}
