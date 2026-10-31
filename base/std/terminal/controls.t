package terminal

enum ControlKind {
    CursorPosition
    CursorUp
    CursorDown
    CursorForward
    CursorBack
    CursorVisible
    ClearScreen
    ClearLine
    ResetStyle
    Bold
    Foreground
    Background
    ForegroundRGB
    BackgroundRGB
    AlternateScreen
    BracketedPaste
    MouseTracking
}

struct Control {
    kind ControlKind
    first u64
    second u64
    third u64
    enabled bool
}

enum Error: error {
    InvalidArgument
    BufferTooSmall
    Closed
    NotTerminal
    Unsupported
    SystemFailure
}

fn encode_control(control Control, buffer u8[]) (u64, error) {
    encoded := zero<u8[64]>()
    encoded[0] = 27
    encoded[1] = 91
    count u64 := 2
    kind := control.kind
    if kind == ControlKind::CursorPosition {
        if control.first == 0 or control.second == 0 or control.first > 65535 or control.second > 65535 {
            return 0, error(Error::InvalidArgument, "cursor coordinates must be in 1..65535")
        }
        count = encode_decimal(control.first, encoded[:], count)
        encoded[count] = 59
        count = encode_decimal(control.second, encoded[:], count + 1)
        encoded[count] = 72
        count = count + 1
    } elif kind == ControlKind::CursorUp or kind == ControlKind::CursorDown or kind == ControlKind::CursorForward or kind == ControlKind::CursorBack {
        if control.first == 0 or control.first > 65535 {
            return 0, error(Error::InvalidArgument, "cursor distance must be in 1..65535")
        }
        count = encode_decimal(control.first, encoded[:], count)
        encoded[count] = 65
        if kind == ControlKind::CursorDown {
            encoded[count] = 66
        } elif kind == ControlKind::CursorForward {
            encoded[count] = 67
        } elif kind == ControlKind::CursorBack {
            encoded[count] = 68
        }
        count = count + 1
    } elif kind == ControlKind::ClearScreen or kind == ControlKind::ClearLine {
        encoded[2] = 50
        encoded[3] = 74
        if kind == ControlKind::ClearLine {
            encoded[3] = 75
        }
        count = 4
    } elif kind == ControlKind::ResetStyle or kind == ControlKind::Bold {
        value u64 := 0
        if kind == ControlKind::Bold {
            value = 22
            if control.enabled {
                value = 1
            }
        }
        count = encode_decimal(value, encoded[:], count)
        encoded[count] = 109
        count = count + 1
    } elif kind == ControlKind::Foreground or kind == ControlKind::Background or kind == ControlKind::ForegroundRGB or kind == ControlKind::BackgroundRGB {
        if control.first > 255 or control.second > 255 or control.third > 255 {
            return 0, error(Error::InvalidArgument, "color components must be in 0..255")
        }
        base u64 := 38
        if kind == ControlKind::Background or kind == ControlKind::BackgroundRGB {
            base = 48
        }
        count = encode_decimal(base, encoded[:], count)
        encoded[count] = 59
        count = count + 1
        rgb := kind == ControlKind::ForegroundRGB or kind == ControlKind::BackgroundRGB
        encoded[count] = 53
        if rgb {
            encoded[count] = 50
        }
        encoded[count + 1] = 59
        count = encode_decimal(control.first, encoded[:], count + 2)
        if rgb {
            encoded[count] = 59
            count = encode_decimal(control.second, encoded[:], count + 1)
            encoded[count] = 59
            count = encode_decimal(control.third, encoded[:], count + 1)
        }
        encoded[count] = 109
        count = count + 1
    } elif kind == ControlKind::CursorVisible or kind == ControlKind::AlternateScreen or kind == ControlKind::BracketedPaste or kind == ControlKind::MouseTracking {
        number u64 := 25
        if kind == ControlKind::AlternateScreen {
            number = 1049
        } elif kind == ControlKind::BracketedPaste {
            number = 2004
        } elif kind == ControlKind::MouseTracking {
            number = 1000
        }
        count = encode_private(number, control.enabled, encoded[:], 0)
        if kind == ControlKind::MouseTracking {
            count = encode_private(1006, control.enabled, encoded[:], count)
        }
    } else {
        return 0, error(Error::InvalidArgument, "unknown terminal control")
    }
    if len(buffer) < count {
        return 0, error(Error::BufferTooSmall, "terminal control buffer too small")
    }
    i u64 := 0
    while i < count {
        buffer[i] = encoded[i]
        i = i + 1
    }
    return count, nil
}

fn encode_private(number u64, enabled bool, buffer u8[], offset u64) u64 {
    buffer[offset] = 27
    buffer[offset + 1] = 91
    buffer[offset + 2] = 63
    end := encode_decimal(number, buffer, offset + 3)
    buffer[end] = 108
    if enabled {
        buffer[end] = 104
    }
    return end + 1
}

fn encode_decimal(value u64, buffer u8[], offset u64) u64 {
    digits u8[10] := .{48, 49, 50, 51, 52, 53, 54, 55, 56, 57}
    reversed := zero<u8[20]>()
    count u64 := 0
    remaining := value
    while count == 0 or remaining != 0 {
        digit := remaining - (remaining / 10) * 10
        reversed[count] = digits[digit]
        count = count + 1
        remaining = remaining / 10
    }
    i u64 := 0
    while i < count {
        buffer[offset + i] = reversed[count - i - 1]
        i = i + 1
    }
    return offset + count
}
