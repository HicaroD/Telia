package boxes

use remote "pkg::other"

type Remote = remote::User

struct Box<T> {
    value T
}

fn helper() string {
    return "library"
}

fn identity<T>(value T) T {
    return value
}

fn label<T>(value T) string {
    return helper()
}

fn remote_box<T>(value Remote) Box<Remote> {
    return Box<Remote>{ value: value }
}
