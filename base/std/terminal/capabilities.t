package terminal

use "std::ffi::c"

@[header="telia_terminal.h", default_cc="c"]
extern {
    @[link_name="telia_terminal_profile"]
    fn native_profile(handle c::c_intptr_t, dumb *bool, ansi *bool, indexed *bool, rgb *bool) i32
    @[link_name="telia_terminal_prepare_output"]
    fn native_prepare_output(handle c::c_intptr_t, saved *rawptr) i32
}

enum Support {
    Unknown
    Unsupported
    Supported
}

struct Capabilities {
    controls Support
    indexed_color Support
    true_color Support
    alternate_screen Support
    bracketed_paste Support
    mouse Support
}

fn capabilities(handle *Handle) (Capabilities, error) {
    result := Capabilities{}
    if handle.closed or not handle.valid {
        return result, error(Error::Closed, "terminal handle is closed")
    }
    dumb := false
    ansi := false
    indexed := false
    rgb := false
    status := native_profile(handle.native, &dumb, &ansi, &indexed, &rgb)
    if status != 0 {
        return result, status_error(status)
    }
    if dumb {
        result.controls = Support::Unsupported
        result.indexed_color = Support::Unsupported
        result.true_color = Support::Unsupported
        result.alternate_screen = Support::Unsupported
        result.bracketed_paste = Support::Unsupported
        result.mouse = Support::Unsupported
    } else {
        if ansi {
            result.controls = Support::Supported
        }
        if indexed {
            result.indexed_color = Support::Supported
        }
        if rgb {
            result.true_color = Support::Supported
        }
    }
    return result, nil
}

fn send_control(handle *Handle, control Control, caps *Capabilities) error {
    attached, err := is_terminal(handle)
    if err != nil {
        return err
    }
    if not attached {
        return error(Error::NotTerminal, "handle is not a terminal")
    }
    support := caps.controls
    if control.kind == ControlKind::Foreground or control.kind == ControlKind::Background {
        support = caps.indexed_color
    } elif control.kind == ControlKind::ForegroundRGB or control.kind == ControlKind::BackgroundRGB {
        support = caps.true_color
    } elif control.kind == ControlKind::AlternateScreen {
        support = caps.alternate_screen
    } elif control.kind == ControlKind::BracketedPaste {
        support = caps.bracketed_paste
    } elif control.kind == ControlKind::MouseTracking {
        support = caps.mouse
    }
    if caps.controls == Support::Unsupported or support == Support::Unsupported {
        return error(Error::Unsupported, "terminal control is unsupported")
    }
    buffer := zero<u8[64]>()
    count, encode_err := encode_control(control, buffer[:])
    if encode_err != nil {
        return encode_err
    }
    bytes := buffer[0:count] @fail
    return write_all(handle, bytes)
}

fn prepare_output(handle *Handle) (SavedMode, error) {
    result := SavedMode{}
    if handle.closed or not handle.valid {
        return result, error(Error::Closed, "terminal handle is closed")
    }
    status := native_prepare_output(handle.native, &result.state)
    return result, status_error(status)
}
