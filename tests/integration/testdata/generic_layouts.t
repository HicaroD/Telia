package main

use "std::io"
use "std::libc"

struct Envelope<T> {
    payload Box<T>
}

struct Box<T> {
    value T
}

struct Node<T> {
    value T
    next *Node<T>
}

type Text = string

fn main() i32 {
    item := Envelope<Text>.{ payload: Box<string>.{ value: "nested" } }
    node := Node<i32>.{ value: 9 }
    io::println(item.payload.value)
    libc::printf("%d\n", node.value)
    return 0
}
