package main

use "std::strings"
use "std::lists"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
}

fn main() i32 {
    owner := strings::clone("  a,a,a,a,a,a,a,a,a  ") @fail
    original := strings::view(&owner)
    _test_fail_after(0)
    trimmed := strings::trim_space(original) @fail
    if trimmed != "a,a,a,a,a,a,a,a,a" or _test_live_allocations() != 1 {
        return 1
    }
    if not strings::contains(trimmed, ",") or not strings::starts_with(trimmed, "a,") or not strings::ends_with(trimmed, ",a") {
        return 2
    }
    invalid, invalid_err := strings::split(trimmed, "")
    if invalid_err == nil or lists::count<string>(&invalid) != 0 or lists::capacity<string>(&invalid) != 0 {
        return 3
    }
    lists::free<string>(&invalid)
    empty, empty_err := strings::split("", "")
    if empty_err == nil or lists::count<string>(&empty) != 0 or lists::capacity<string>(&empty) != 0 {
        return 10
    }
    lists::free<string>(&empty)
    if _test_live_allocations() != 1 {
        return 4
    }
    budget i64 := 0
    saw_failure := false
    saw_success := false
    while budget <= 8 {
        _test_fail_after(budget)
        fields, err := strings::split(trimmed, ",")
        if err != nil {
            saw_failure = true
            if lists::count<string>(&fields) != 0 or lists::capacity<string>(&fields) != 0 or _test_live_allocations() != 1 {
                return 5
            }
        } else {
            saw_success = true
            if lists::count<string>(&fields) != 9 or _test_live_allocations() != 2 {
                return 6
            }
            saved := lists::get<string>(&fields, 0) @fail
            last := lists::get<string>(&fields, 8) @fail
            lists::free<string>(&fields)
            if saved != "a" or last != "a" or _test_live_allocations() != 1 {
                return 7
            }
        }
        lists::free<string>(&fields)
        if strings::view(&owner) != "  a,a,a,a,a,a,a,a,a  " or trimmed != "a,a,a,a,a,a,a,a,a" {
            return 8
        }
        budget = budget + 1
    }
    _test_fail_after(-1)
    strings::free(&owner)
    if not saw_failure or not saw_success or _test_live_allocations() != 0 {
        return 9
    }
    return 0
}
