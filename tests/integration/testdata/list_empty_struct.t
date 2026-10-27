package main

use "std::lists"
use "std::libc"

struct Marker {
}

fn main() i32 {
    values := lists::create<Marker>()
    defer lists::free<Marker>(&values)
    lists::append<Marker>(&values, Marker.{}) @fail
    libc::printf("%llu\n%llu\n", lists::count<Marker>(&values), sizeof<Marker>())
    return 0
}
