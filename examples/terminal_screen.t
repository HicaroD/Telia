package main

use "std::terminal"

fn draw(output *terminal::Handle, caps *terminal::Capabilities) error {
    terminal::send_control(output, terminal::Control{kind: terminal::ControlKind::ClearScreen}, caps) @catch err {
        return err
    }
    terminal::send_control(output, terminal::Control{kind: terminal::ControlKind::CursorPosition, first: 1, second: 1}, caps) @catch err {
        return err
    }
    return terminal::write_text(output, "Telia terminal demo — resize to redraw; press q or Ctrl+C to quit.")
}

fn run() error {
    input := terminal::open_controlling(true) @catch err {
        return err
    }
    defer terminal::close(&input)
    output := terminal::open_controlling(false) @catch err {
        return err
    }
    defer terminal::close(&output)
    raw := terminal::apply_mode(&input, terminal::raw_mode()) @catch err {
        return err
    }
    defer terminal::restore(&raw)
    prepared := terminal::prepare_output(&output) @catch err {
        return err
    }
    defer terminal::restore(&prepared)
    caps := terminal::capabilities(&output) @catch err {
        return err
    }

    terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::AlternateScreen, enabled: true}, &caps) @catch err {
        return err
    }
    defer terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::AlternateScreen}, &caps)
    terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::CursorVisible}, &caps) @catch err {
        return err
    }
    defer terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::CursorVisible, enabled: true}, &caps)

    draw(&output, &caps) @catch err {
        return err
    }
    buffer := zero<u8[256]>()
    while true {
        event := terminal::wait(&input, terminal::WaitOptions{input: true, resize: true, timeout_ms: -1}) @catch err {
            return err
        }
        if event.resized {
            draw(&output, &caps) @catch err {
                return err
            }
        }
        if event.input {
            result := terminal::read(&input, buffer[:]) @catch err {
                return err
            }
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
