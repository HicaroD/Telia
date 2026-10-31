package main

use "std::terminal"

fn main() i32 {
    input := terminal::stdin()
    output := terminal::stdout()
    buffer := zero<u8[4096]>()
    while true {
        result := terminal::read(&input, buffer[:]) @fail
        if result.eof {
            return 0
        }
        bytes := buffer[0:result.count] @fail
        terminal::write_all(&output, bytes) @fail
    }
    return 0
}
