package main

use "std::io"

struct User {
  name string
  age  int
}

fn main() {
  user := User.{
    name: "Hicaro",
    age: 21,
  }
  io::println(user.name)
  return
}
