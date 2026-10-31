package math

use "std::testing"

fn check_sum(a i32, b i32, expected i32) {
    testing::assert_eq(add(a, b), expected)
}

test "adds positive and negative numbers" {
    check_sum(2, 3, 5)
    check_sum(-2, 3, 1)
}

test "supports boolean and inequality assertions" {
    testing::assert(add(2, 3) > 0)
    testing::assert_ne(add(2, 3), 0)
    testing::assert_eq("Telia", "Telia")
}
