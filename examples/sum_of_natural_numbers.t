pkg main;

#[default_cc="c"]
extern C {
  fn puts(s *u8);
  fn printf(format *u8, ...);
}

fn sum_of_natural_numbers(n int) int {
  sum, i := 0, 1;

  while i <= n {
    sum = sum + i;
    i = i + 1;
  }

  return sum;
}

fn main() {
  n, sum := 10, sum_of_natural_numbers(n);
  C.printf("Sum of %d first natural numbers: %d", n, sum);
  return;
}
