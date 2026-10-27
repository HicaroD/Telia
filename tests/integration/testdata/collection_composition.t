package main

use "std::lists"
use "std::io"
use "std::libc"

struct Record {
    name string
    values i32[2]
}

fn main() i32 {
    record := Record.{name: "record", values: .{1, 2}}
    records := lists::create<Record>()
    defer lists::free<Record>(&records)
    lists::append<Record>(&records, record) @fail
    record.values[1] = 9
    view := lists::view<Record>(&records)
    for item in view {
        io::println(item.name)
        libc::printf("%d\n", item.values[1])
    }
    arrays := lists::create<i32[2]>()
    defer lists::free<i32[2]>(&arrays)
    lists::append<i32[2]>(&arrays, record.values) @fail
    array := lists::get<i32[2]>(&arrays, 0) @fail
    libc::printf("%d\n", array[1])
    inner := lists::create<string>()
    defer lists::free<string>(&inner)
    lists::append<string>(&inner, "nested") @fail
    outer := lists::create<lists::List<string>>()
    defer lists::free<lists::List<string>>(&outer)
    lists::append<lists::List<string>>(&outer, inner) @fail
    stored := lists::get<lists::List<string>>(&outer, 0) @fail
    text := lists::get<string>(&stored, 0) @fail
    io::println(text)
    pointers := lists::create<*Record>()
    defer lists::free<*Record>(&pointers)
    lists::append<*Record>(&pointers, &record) @fail
    pointer := lists::get<*Record>(&pointers, 0) @fail
    io::println(pointer.name)
    return 0
}
