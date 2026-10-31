package process

use "std::strings"
use "std::errors"
use "std::unicode::utf8"

struct SpawnOptions {
    pipe_stdin bool
    pipe_stdout bool
    pipe_stderr bool
}

// An owning process handle. Do not copy a live Child.
struct Child {
    state rawptr
}

struct ExitStatus {
    code u32
    signal u32
}

fn validate_argument(value string) error {
    if strings::contains(value, "\0") {
        return error(errors::Error::EmbeddedZero, "embedded zero in process argument")
    }
    return utf8::validate(value)
}

fn spawn(executable string, arguments string[], options SpawnOptions) (Child, error) {
    child := Child{}
    if strings::len(executable) == 0 {
        return child, error(errors::Error::InvalidPath, "empty executable")
    }
    err := validate_argument(executable)
    index u64 := 0
    while err == nil and index < len(arguments) {
        err = validate_argument(arguments[index])
        index = index + 1
    }
    if err == nil {
        err = spawn_into(executable, slice_data<string>(arguments), len(arguments), options.pipe_stdin, options.pipe_stdout, options.pipe_stderr, &child.state)
    }
    return child, err
}

fn wait(child *Child) (ExitStatus, error) {
    result := ExitStatus{}
    done := false
    err := wait_into(child.state, true, &done, &result.code, &result.signal)
    return result, err
}

fn poll(child *Child) (ExitStatus, bool, error) {
    result := ExitStatus{}
    done := false
    err := wait_into(child.state, false, &done, &result.code, &result.signal)
    return result, done, err
}

fn close_child(child *Child) error {
    return dispose_into(&child.state)
}

fn read_stdout(child *Child, buffer u8[]) (u64, bool, bool, error) {
    return read_pipe(child, 1, buffer)
}

fn read_stderr(child *Child, buffer u8[]) (u64, bool, bool, error) {
    return read_pipe(child, 2, buffer)
}

fn read_pipe(child *Child, stream u32, buffer u8[]) (u64, bool, bool, error) {
    count u64 := 0
    eof := false
    would_block := false
    err := pipe_read(child.state, stream, slice_data<u8>(buffer), len(buffer), &count, &eof, &would_block)
    return count, eof, would_block, err
}

fn write_stdin(child *Child, data string) (u64, bool, error) {
    count u64 := 0
    would_block := false
    err := pipe_write(child.state, data, &count, &would_block)
    return count, would_block, err
}

fn close_stdin(child *Child) error {
    return pipe_close(child.state, 0)
}

fn terminate(child *Child) error {
    return terminate_into(child.state, false)
}

fn kill(child *Child) error {
    return terminate_into(child.state, true)
}

@[default_cc="c"]
extern {
    @[link_name="_telia_child_spawn"]
    fn spawn_into(executable string, arguments rawptr, count u64, input bool, output bool, stderr bool, result *rawptr) error
    @[link_name="_telia_child_wait"]
    fn wait_into(state rawptr, blocking bool, done *bool, code *u32, signal *u32) error
    @[link_name="_telia_child_dispose"]
    fn dispose_into(state *rawptr) error
    @[link_name="_telia_child_read"]
    fn pipe_read(state rawptr, stream u32, buffer rawptr, size u64, count *u64, eof *bool, would_block *bool) error
    @[link_name="_telia_child_write"]
    fn pipe_write(state rawptr, data string, count *u64, would_block *bool) error
    @[link_name="_telia_child_close_pipe"]
    fn pipe_close(state rawptr, stream u32) error
    @[link_name="_telia_child_terminate"]
    fn terminate_into(state rawptr, force bool) error
}
