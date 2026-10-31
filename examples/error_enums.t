package main

use "std::io"

enum FileError: error {
    NotFound {
        path string
    }
    PermissionDenied
}

fn read_file(path string) (string, error) {
    return "", error(FileError::NotFound{path: path}, "file does not exist")
}

fn main() i32 {
    text := read_file("settings.txt") @catch err {
        description := match err {
            FileError::NotFound{path} => path
            FileError::PermissionDenied => "permission denied"
            _ => err.msg
        }
        io::println(description)
        io::println(err.msg)
        return 0
    }
    io::println(text)
    return 0
}
