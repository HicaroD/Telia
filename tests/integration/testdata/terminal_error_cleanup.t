package main

use "std::terminal"

fn run(input *terminal::Handle, output *terminal::Handle) error {
    saved := terminal::apply_mode(input, terminal::raw_mode()) @try
    defer terminal::restore(&saved)
    terminal::write_text(output, "R") @try
    caps := terminal::Capabilities{}
    terminal::send_control(output, terminal::Control{kind: terminal::ControlKind::ClearScreen}, &caps) @try
    return nil
}

fn main() i32 {
    input := terminal::stdin()
    output := terminal::stdout()
    err := run(&input, &output)
    if err == nil {
        return 1
    }
    return 0
}
