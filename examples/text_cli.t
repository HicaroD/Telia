package main

use "std::io"
use "std::strings"
use "std::lists"
use "std::process"
use "std::fs"
use "std::strconv"
use "std::unicode::utf8"

fn report(message string) i32 {
    diagnostic_err := io::write_stderr(message)
    // A failed diagnostic is still an unsuccessful invocation; never @fail.
    return 1
}

fn print_line(text string, repeat i64, counts bool) i32 {
    formatted := zero<strings::Owned>()
    defer strings::free(&formatted)
    output := text
    if counts {
        count, count_err := utf8::count(text)
        if count_err != nil {
            return report("invalid UTF-8 input\n")
        }
        value, format_err := strconv::format_u64(count)
        formatted = value
        if format_err != nil {
            return report("could not format count\n")
        }
        output = strings::view(&formatted)
    }
    iteration i64 := 0
    while iteration < repeat {
        if io::write_stdout(output) != nil {
            return report("stdout write failed\n")
        }
        if io::write_stdout("\n") != nil {
            return report("stdout write failed\n")
        }
        iteration = iteration + 1
    }
    return 0
}

fn print_lines(text string, repeat i64, counts bool) i32 {
    if utf8::validate(text) != nil {
        return report("invalid UTF-8 input\n")
    }
    lines, split_err := strings::split(text, "\n")
    if split_err != nil {
        return report("could not split input\n")
    }
    defer lists::free<string>(&lines)
    index u64 := 0
    while index < lists::count<string>(&lines) {
        line, get_err := lists::get<string>(&lines, index)
        if get_err != nil {
            return report("could not read line\n")
        }
        trimmed, trim_err := strings::trim_space(line)
        if trim_err != nil {
            return report("invalid UTF-8 input\n")
        }
        status := print_line(trimmed, repeat, counts)
        if status != 0 {
            return status
        }
        index = index + 1
    }
    return 0
}

fn main() i32 {
    args, args_err := process::args()
    if args_err != nil {
        return report("could not read arguments\n")
    }
    defer process::free(&args)
    if process::count(&args) > 3 {
        return report("usage: text_cli [INPUT|- [REPEAT [COUNTS]]]\n")
    }
    path := "-"
    if process::count(&args) > 0 {
        value, access_err := process::get(&args, 0)
        if access_err != nil {
            return report("could not read arguments\n")
        }
        path = value
    }
    default_repeat := zero<strings::Owned>()
    defer strings::free(&default_repeat)
    repeat_text := "1"
    if process::count(&args) > 1 {
        value, access_err := process::get(&args, 1)
        if access_err != nil {
            return report("could not read arguments\n")
        }
        repeat_text = value
    } else {
        value, found, lookup_err := process::lookup_env("TELIA_TEXT_REPEAT")
        // Explicitly transfer this owner; do not free its shallow copy.
        default_repeat = value
        if lookup_err != nil {
            return report("could not read environment\n")
        }
        if found {
            repeat_text = strings::view(&default_repeat)
        }
    }
    repeat, repeat_err := strconv::parse_i64(repeat_text)
    if repeat_err != nil or repeat < 1 or repeat > 100 {
        return report("repeat must be an integer from 1 to 100\n")
    }
    counts := false
    if process::count(&args) > 2 {
        value, access_err := process::get(&args, 2)
        if access_err != nil {
            return report("could not read arguments\n")
        }
        parsed, parse_err := strconv::parse_bool(value)
        if parse_err != nil {
            return report("counts must be true or false\n")
        }
        counts = parsed
    }
    input, read_err := read_input(path)
    if read_err != nil {
        return report("input read failed\n")
    }
    defer strings::free(&input)
    return print_lines(strings::view(&input), repeat, counts)
}

fn read_input(path string) (strings::Owned, error) {
    if path == "-" {
        return io::read_stdin()
    }
    return fs::read_file(path)
}
