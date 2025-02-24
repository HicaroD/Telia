package main;

#[default_cc="c"]
extern C {
  fn puts(s cstring);
  fn printf(format cstring, ...);
}

fn sum_of_natural_numbers(n int) int {
  sum, i := 0, 1;

  while i <= n {
    sum = sum + i;
    i = i + 1;
  }

  return sum;
}

fn formula(n int) int {
  return (n * (n + 1))/2;
}

fn main() {
  n, sum := 10, sum_of_natural_numbers(n);
  C.printf("Sum of %d first natural numbers: %d", n, sum);

  sum = formula(n);
  C.printf("Sum of %d first natural numbers with formula: %d", n, sum);
  return;
}
