package terminal

use "std::ffi::c"

@[header="telia_terminal.h", default_cc="c"]
extern {
    @[link_name="telia_terminal_is_terminal"]
    fn native_is_terminal(handle c::c_intptr_t, result *bool) i32
    @[link_name="telia_terminal_size"]
    fn native_size(handle c::c_intptr_t, columns *u64, rows *u64) i32
    @[link_name="telia_terminal_get_mode"]
    fn native_get_mode(handle c::c_intptr_t, echo *bool, canonical *bool, signals *bool) i32
    @[link_name="telia_terminal_apply_mode"]
    fn native_apply_mode(handle c::c_intptr_t, raw bool, echo bool, canonical bool, signals bool, saved *rawptr) i32
    @[link_name="telia_terminal_restore"]
    fn native_restore(saved *rawptr) i32
}

struct Mode {
    echo bool
    canonical bool
    signals bool
    raw bool
}

struct SavedMode {
    state rawptr
}

struct Size {
    columns u64
    rows u64
}

fn is_terminal(handle *Handle) (bool, error) {
    if handle.closed or not handle.valid {
        return false, error(Error::Closed, "terminal handle is closed")
    }
    result := false
    status := native_is_terminal(handle.native, &result)
    return result, status_error(status)
}

fn get_mode(handle *Handle) (Mode, error) {
    result := Mode{}
    if handle.closed or not handle.valid {
        return result, error(Error::Closed, "terminal handle is closed")
    }
    status := native_get_mode(handle.native, &result.echo, &result.canonical, &result.signals)
    return result, status_error(status)
}

fn raw_mode() Mode {
    return Mode{raw: true}
}

fn apply_mode(handle *Handle, mode Mode) (SavedMode, error) {
    result := SavedMode{}
    if handle.closed or not handle.valid {
        return result, error(Error::Closed, "terminal handle is closed")
    }
    status := native_apply_mode(handle.native, mode.raw, mode.echo, mode.canonical, mode.signals, &result.state)
    return result, status_error(status)
}

fn restore(saved *SavedMode) error {
    return status_error(native_restore(&saved.state))
}

fn size(handle *Handle) (Size, error) {
    result := Size{}
    if handle.closed or not handle.valid {
        return result, error(Error::Closed, "terminal handle is closed")
    }
    status := native_size(handle.native, &result.columns, &result.rows)
    return result, status_error(status)
}

fn status_error(status i32) error {
    if status == 0 {
        return nil
    }
    if status == 1 {
        return error(Error::InvalidArgument, "invalid terminal argument")
    }
    if status == 4 {
        return error(Error::NotTerminal, "handle is not a terminal")
    }
    if status == 5 {
        return error(Error::Unsupported, "terminal operation is unsupported")
    }
    return error(Error::SystemFailure, "terminal system operation failed")
}
