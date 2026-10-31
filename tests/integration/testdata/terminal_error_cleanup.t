package main

use "std::terminal"

fn run(input *terminal::Handle, output *terminal::Handle) error {
    saved := terminal::apply_mode(input, terminal::raw_mode()) @catch err {
        return err
    }
    defer terminal::restore(&saved)
    terminal::write_text(output, "R") @catch err {
        return err
    }
    caps := terminal::Capabilities{}
    return terminal::send_control(output, terminal::Control{kind: terminal::ControlKind::ClearScreen}, &caps)
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
