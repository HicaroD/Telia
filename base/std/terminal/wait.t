package terminal

use "std::ffi::c"

@[header="telia_terminal.h", default_cc="c"]
extern {
    @[link_name="telia_terminal_wait"]
    fn native_wait(handle c::c_intptr_t, input bool, resize bool, timeout_ms i64, buffered_input bool, initialized *bool, columns *u64, rows *u64, ready *bool, resized *bool, timed_out *bool) i32
}

struct WaitOptions {
    input bool
    resize bool
    timeout_ms i64
}

struct WaitResult {
    input bool
    resized bool
    timed_out bool
}

fn wait(handle *Handle, options WaitOptions) (WaitResult, error) {
    result := WaitResult{}
    if handle.closed or not handle.valid {
        return result, error(Error::Closed, "terminal handle is closed")
    }
    buffered_unit := handle.pending_utf16 != 0 and (handle.pending_utf16 < 55297 or handle.pending_utf16 > 56320)
    buffered := handle.pending_utf8_count != 0 or buffered_unit
    status := native_wait(handle.native, options.input, options.resize, options.timeout_ms, buffered, &handle.resize_initialized, &handle.last_columns, &handle.last_rows, &result.input, &result.resized, &result.timed_out)
    return result, status_error(status)
}
