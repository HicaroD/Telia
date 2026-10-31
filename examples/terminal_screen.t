package main

use "std::terminal"

fn draw(output *terminal::Handle, caps *terminal::Capabilities) error {
    terminal::send_control(output, terminal::Control{kind: terminal::ControlKind::ClearScreen}, caps) @try
    terminal::send_control(output, terminal::Control{kind: terminal::ControlKind::CursorPosition, first: 1, second: 1}, caps) @try
    return terminal::write_text(output, "Telia terminal demo — resize to redraw; press q or Ctrl+C to quit.")
}

fn run() error {
    input := terminal::open_controlling(true) @try
    defer terminal::close(&input)
    output := terminal::open_controlling(false) @try
    defer terminal::close(&output)
    raw := terminal::apply_mode(&input, terminal::raw_mode()) @try
    defer terminal::restore(&raw)
    prepared := terminal::prepare_output(&output) @try
    defer terminal::restore(&prepared)
    caps := terminal::capabilities(&output) @try

    terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::AlternateScreen, enabled: true}, &caps) @try
    defer terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::AlternateScreen}, &caps)
    terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::CursorVisible}, &caps) @try
    defer terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::CursorVisible, enabled: true}, &caps)

    draw(&output, &caps) @try
    buffer := zero<u8[256]>()
    while true {
        event := terminal::wait(&input, terminal::WaitOptions{input: true, resize: true, timeout_ms: -1}) @try
        if event.resized {
            draw(&output, &caps) @try
        }
        if event.input {
            result := terminal::read(&input, buffer[:]) @try
            if result.eof {
                return nil
            }
            i u64 := 0
            while i < result.count {
                if buffer[i] == 113 or buffer[i] == 3 {
                    return nil
                }
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
