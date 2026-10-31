package main

use "std::terminal"

fn main() i32 {
    output := terminal::stdout()
    attached := terminal::is_terminal(&output) @fail
    if attached {
        return 1
    }
    caps := terminal::capabilities(&output) @fail
    if caps.controls != terminal::Support::Unsupported or caps.true_color != terminal::Support::Unsupported {
        return 2
    }
    caps.controls = terminal::Support::Supported
    err := terminal::send_control(&output, terminal::Control{kind: terminal::ControlKind::ClearScreen}, &caps)
    if err == nil {
        return 3
    }
    dimensions, size_err := terminal::size(&output)
    if size_err == nil {
        return 4
    }
    saved, mode_err := terminal::apply_mode(&output, terminal::raw_mode())
    if mode_err == nil {
        return 5
    }
    prepared, prepare_err := terminal::prepare_output(&output)
    if prepare_err == nil {
        return 6
    }
    invalid := terminal::Handle{}
    buffer := zero<u8[1]>()
    count, write_err := terminal::write(&invalid, buffer[:])
    if write_err == nil or count != 0 {
        return 7
    }
    wait, wait_err := terminal::wait(&output, terminal::WaitOptions{input: true, timeout_ms: -2})
    if wait_err == nil {
        return 8
    }
    return 0
}
