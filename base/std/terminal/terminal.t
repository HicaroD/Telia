package terminal

use "std::ffi::c"
use "std::strings"

struct Handle {
    native c::c_intptr_t
    owned bool
    closed bool
    valid bool
    resize_initialized bool
    last_columns u64
    last_rows u64
}

@[header="telia_terminal.h", default_cc="c"]
extern {
    @[link_name="telia_terminal_standard"]
    fn native_standard(stream i32) c::c_intptr_t

    @[link_name="telia_terminal_write"]
    fn native_write(handle c::c_intptr_t, bytes rawptr, length u64, count *u64) i32

    @[link_name="telia_terminal_write_text"]
    fn native_write_text(handle c::c_intptr_t, bytes c::CString, length u64, count *u64) i32

    @[link_name="telia_terminal_read"]
    fn native_read(handle c::c_intptr_t, bytes rawptr, length u64, count *u64, eof *bool) i32

    @[link_name="telia_terminal_close"]
    fn native_close(handle c::c_intptr_t) i32

    @[link_name="telia_terminal_open"]
    fn native_open(input bool, handle *c::c_intptr_t) i32
}

fn stdout() Handle {
    return borrow_native(native_standard(1))
}

fn stdin() Handle {
    return borrow_native(native_standard(0))
}

fn stderr() Handle {
    return borrow_native(native_standard(2))
}

fn borrow_native(value c::c_intptr_t) Handle {
    return Handle{native: value, valid: true}
}

fn close(handle *Handle) error {
    if handle.closed or not handle.valid {
        return nil
    }
    handle.closed = true
    if handle.owned {
        if native_close(handle.native) != 0 {
            return error(Error::SystemFailure, "terminal close failed")
        }
    }
    return nil
}

fn open_controlling(input bool) (Handle, error) {
    result := Handle{}
    status := native_open(input, &result.native)
    if status != 0 {
        return result, status_error(status)
    }
    result.valid = true
    result.owned = true
    return result, nil
}

struct ReadResult {
    count u64
    eof bool
}

fn read(handle *Handle, bytes u8[]) (ReadResult, error) {
    result := ReadResult{}
    if handle.closed or not handle.valid {
        return result, error(Error::Closed, "terminal handle is closed")
    }
    status := native_read(handle.native, slice_data<u8>(bytes), len(bytes), &result.count, &result.eof)
    if status != 0 {
        return result, error(Error::SystemFailure, "terminal read failed")
    }
    return result, nil
}

fn write(handle *Handle, bytes u8[]) (u64, error) {
    if handle.closed or not handle.valid {
        return 0, error(Error::Closed, "terminal handle is closed")
    }
    count u64 := 0
    status := native_write(handle.native, slice_data<u8>(bytes), len(bytes), &count)
    if status != 0 {
        return count, error(Error::SystemFailure, "terminal write failed")
    }
    return count, nil
}

fn write_all(handle *Handle, bytes u8[]) error {
    if handle.closed or not handle.valid {
        return error(Error::Closed, "terminal handle is closed")
    }
    offset u64 := 0
    while offset < len(bytes) {
        remaining := bytes[offset:] @fail
        count, err := write(handle, remaining)
        if err != nil {
            return err
        }
        if count == 0 {
            return error(Error::SystemFailure, "terminal write made no progress")
        }
        offset = offset + count
    }
    return nil
}

fn write_text(handle *Handle, text string) error {
    if handle.closed or not handle.valid {
        return error(Error::Closed, "terminal handle is closed")
    }
    offset u64 := 0
    length := strings::len(text)
    while offset < length {
        remaining := strings::slice(text, offset, length) @fail
        count u64 := 0
        status := native_write_text(handle.native, strings::c_pointer(remaining), length - offset, &count)
        if status != 0 {
            return status_error(status)
        }
        if count == 0 {
            return error(Error::SystemFailure, "terminal write made no progress")
        }
        offset = offset + count
    }
    return nil
}
