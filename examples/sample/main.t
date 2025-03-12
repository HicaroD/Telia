package main

use "std::io"
use "pkg::utils"

fn main() i32 {
  io::println("Hello, world 😃")
  from_main_package()
  utils::foo()
  user := utils::User{
    name: "Hicaro"
  }
  io::println(user.name)
  return 0
}
