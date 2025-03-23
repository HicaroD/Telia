package main

use "std::io"

struct User {
  name string
  age  int
}

fn main() {
  user := User{
    name: "Hicaro",
    age: 21,
  }

  username := &user.name
  *username = "Hicro Danrlley"
  io::puts(user.name)
  return
}
