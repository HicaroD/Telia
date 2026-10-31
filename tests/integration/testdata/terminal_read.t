package main

use "std::terminal"

fn main() i32 {
    input := terminal::stdin()
    output := terminal::stdout()
    buffer := zero<u8[2]>()
    empty := buffer[0:0] @fail
    initial := terminal::read(&input, empty) @fail
    if initial.count != 0 or initial.eof {
        return 1
    }
    while true {
        result := terminal::read(&input, buffer[:]) @fail
        if result.eof {
            return 0
        }
        chunk := buffer[0:result.count] @fail
        terminal::write_all(&output, chunk) @fail
    }
    return 2
}
