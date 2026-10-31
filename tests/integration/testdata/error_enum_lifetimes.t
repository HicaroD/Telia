package main

use "std::strings"

enum Failure: error {
    Small {
        code i32
    }
    Large {
        data i64[256]
        text string
    }
}

fn erase(data *i64[256]) {
    values := *data
    values[255] = 0
    *data = values
    return
}

fn classify<T>(value T, message string) error {
    return error(value, message)
}

fn construct(text string) error {
    data i64[256] := .{}
    data[255] = 99
    defer erase(&data)
    return classify<Failure>(Failure::Large{data: data, text: text}, text)
}

fn forward(err error) (i32, error) {
    return 0, err
}

fn inspect(err error) i64 {
    result := forward(err) @catch caught {
        copy := caught
        value i64 := match copy {
            Failure::Small{code} => 0
            Failure::Large{data, text} => {
                defer erase(&data)
                if text != "borrowed" or caught.msg != text {
                    return 2
                }
                yield data[255]
            }
            _ => 0
        }
        return value
    }
    return 0
}

fn main() i32 {
    owner := strings::clone("borrowed") @fail
    defer strings::free(&owner)
    err := construct(strings::view(&owner))
    if inspect(err) != 99 {
        return 1
    }
    return 0
}
