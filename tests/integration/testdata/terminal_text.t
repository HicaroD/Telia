package main

use "std::terminal"

fn main() i32 {
    output := terminal::stdout()
    terminal::write_text(&output, "A\0é\n") @fail
    return 0
}
