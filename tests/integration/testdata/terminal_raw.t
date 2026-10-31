package main

use "std::terminal"

fn main() i32 {
    input := terminal::stdin()
    output := terminal::stdout()
    attached := terminal::is_terminal(&input) @fail
    if not attached {
        return 1
    }
    dimensions := terminal::size(&input) @fail
    if dimensions.columns != 80 or dimensions.rows != 24 {
        return 2
    }
    before := terminal::get_mode(&input) @fail
    saved := terminal::apply_mode(&input, terminal::raw_mode()) @fail
    defer terminal::restore(&saved)
    during := terminal::get_mode(&input) @fail
    if during.echo or during.canonical or during.signals {
        return 3
    }
    ready u8[1] := .{82}
    terminal::write_all(&output, ready[:]) @fail
    buffer := zero<u8[1]>()
    result := terminal::read(&input, buffer[:]) @fail
    if result.count != 1 or result.eof or buffer[0] != 3 {
        return 4
    }
    terminal::restore(&saved) @fail
    after := terminal::get_mode(&input) @fail
    if before.echo != after.echo or before.canonical != after.canonical or before.signals != after.signals {
        return 5
    }
    return 0
}
