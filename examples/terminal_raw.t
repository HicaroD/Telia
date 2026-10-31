package main

use "std::terminal"
use "std::libc"

fn run() error {
    input := terminal::open_controlling(true) @try
    defer terminal::close(&input)
    saved := terminal::apply_mode(&input, terminal::raw_mode()) @try
    defer terminal::restore(&saved)

    output := terminal::stdout()
    terminal::write_text(&output, "Input bytes (Ctrl+C exits):\n") @try
    buffer := zero<u8[256]>()
    while true {
        ready := terminal::wait(&input, terminal::WaitOptions{input: true, timeout_ms: -1}) @try
        if ready.input {
            result := terminal::read(&input, buffer[:]) @try
            if result.eof {
                return nil
            }
            i u64 := 0
            while i < result.count {
                if buffer[i] == 3 {
                    return nil
                }
                libc::printf("%d\n", buffer[i])
                i = i + 1
            }
        }
    }
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
