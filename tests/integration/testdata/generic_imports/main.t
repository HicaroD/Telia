package main

use "std::io"
use "std::libc"
use "pkg::boxes"
use "pkg::other"

struct User {
    name string
}

fn helper() string {
    return "caller helper"
}

fn main() i32 {
    user := User.{ name: "caller" }
    box boxes::Box<User> := boxes::Box<User>.{ value: user }
    copy := boxes::identity<boxes::Box<User>>(box)
    io::println(copy.value.name)
    io::println(boxes::label<i32>(1))
    io::println(other::label<i32>(1))
    remote := other::User.{ id: 9 }
    remote_box := boxes::Box<other::User>.{ value: remote }
    libc::printf("%d\n", remote_box.value.id)
    alias_copy := boxes::identity<boxes::Remote>(remote)
    libc::printf("%d\n", alias_copy.id)
    return 0
}
