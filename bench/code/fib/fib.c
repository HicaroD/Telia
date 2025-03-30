#include <stdio.h>

int fib(int n) {
  if (n <= 1) {
    return n;
   } else {
      return fib(n-1) + fib(n-2);
   }
}

int main() {
   int n = 35;
   int r = fib(n);
   printf("%d: %d\n", n, r);
}

