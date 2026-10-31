package main

use "std::errors"
use "std::memory"

enum Failure: error {
    Large {
        data i64[256]
    }
}

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
}

fn main() i32 {
    _test_fail_after(0)
    data i64[256] := .{}
    data[255] = 42
    err := error(Failure::Large{data: data}, "")
    copy := err
    value i64 := match copy {
        Failure::Large{data} => data[255]
        _ => 0
    }
    if value != 42 or err == nil {
        return 1
    }
    ptr, allocErr := memory::alloc(16)
    if allocErr.kind != errors::Error::AllocationFailed {
        return 2
    }
    if _test_live_allocations() != 0 {
        return 3
    }
    return 0
}
