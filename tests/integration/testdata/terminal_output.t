package main

use "std::terminal"

fn main() i32 {
    output := terminal::stdout()
    caps := terminal::capabilities(&output) @fail
    if caps.controls != terminal::Support::Supported or caps.indexed_color != terminal::Support::Supported or caps.true_color != terminal::Support::Supported {
        return 1
    }
    saved := terminal::prepare_output(&output) @fail
    defer terminal::restore(&saved)
    caps.true_color = terminal::Support::Unsupported
    rgb := terminal::Control{kind: terminal::ControlKind::ForegroundRGB, first: 12, second: 34, third: 56}
    rejected := terminal::send_control(&output, rgb, &caps)
    if rejected == nil {
        return 2
    }
    caps.true_color = terminal::Support::Supported
    terminal::send_control(&output, rgb, &caps) @fail
    invalid := terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::CursorPosition}, &caps)
    if invalid == nil {
        return 3
    }
    terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::ResetStyle}, &caps) @fail
    return 0
}
