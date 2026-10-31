package main

use "std::strconv"
use "std::strings"
use "std::errors"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
}

fn main() i32 {
    _test_fail_after(0)
    unsigned := strconv::parse_u64("18446744073709551615") @fail
    signed := strconv::parse_i64("-9223372036854775808") @fail
    flag := strconv::parse_bool("true") @fail
    if unsigned != 18446744073709551615 or signed != -9223372036854775807 - 1 or not flag {
        return 1
    }
    budget i64 := 0
    while budget <= 1 {
        _test_fail_after(budget)
        a, ea := strconv::format_i64(-9223372036854775807 - 1)
        if budget == 0 {
            if ea.kind != errors::Error::AllocationFailed or strings::view(&a) != "" or _test_live_allocations() != 0 {
                return 2
            }
        } else {
            if ea != nil or strings::view(&a) != "-9223372036854775808" or _test_live_allocations() != 1 {
                return 3
            }
        }
        strings::free(&a)
        strings::free(&a)
        if strings::view(&a) != "" or _test_live_allocations() != 0 {
            return 4
        }
        _test_fail_after(budget)
        b, eb := strconv::format_u64(18446744073709551615)
        if budget == 0 {
            if eb.kind != errors::Error::AllocationFailed or strings::view(&b) != "" or _test_live_allocations() != 0 {
                return 5
            }
        } else {
            if eb != nil or strings::view(&b) != "18446744073709551615" or _test_live_allocations() != 1 {
                return 6
            }
        }
        strings::free(&b)
        strings::free(&b)
        _test_fail_after(budget)
        c, ec := strconv::format_bool(true)
        if budget == 0 {
            if ec.kind != errors::Error::AllocationFailed or strings::view(&c) != "" or _test_live_allocations() != 0 {
                return 7
            }
        } else {
            if ec != nil or strings::view(&c) != "true" or _test_live_allocations() != 1 {
                return 8
            }
        }
        strings::free(&c)
        strings::free(&c)
        _test_fail_after(budget)
        d, ed := strconv::format_bool(false)
        if budget == 0 {
            if ed.kind != errors::Error::AllocationFailed or strings::view(&d) != "" or _test_live_allocations() != 0 {
                return 9
            }
        } else {
            if ed != nil or strings::view(&d) != "false" or _test_live_allocations() != 1 {
                return 10
            }
        }
        strings::free(&d)
        strings::free(&d)
        if strings::view(&b) != "" or strings::view(&c) != "" or strings::view(&d) != "" or _test_live_allocations() != 0 {
            return 11
        }
        budget = budget + 1
    }
    _test_fail_after(-1)
    first := strconv::format_u64(42) @fail
    second := strconv::format_u64(7) @fail
    strings::free(&second)
    if strings::view(&first) != "42" or _test_live_allocations() != 1 {
        return 12
    }
    strings::free(&first)
    return 0
}
