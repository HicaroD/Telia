package fs

use "std::strings"

// Owning handle: do not copy a live File. close clears it and is repeatable.
struct File {
    state rawptr
}

fn open_read(path string) (File, error) {
    result := File{}
    err := validate_path(path)
    if err == nil {
        err = file_open(path, &result.state)
    }
    return result, err
}

fn read(file *File, buffer u8[]) (u64, bool, error) {
    count u64 := 0
    eof := false
    err := file_read(file.state, slice_data<u8>(buffer), len(buffer), &count, &eof)
    return count, eof, err
}

fn close(file *File) error {
    return file_close(&file.state)
}

// Reads at most limit bytes. No lookahead or whole-file allocation occurs.
fn read_preview(path string, limit u64) (strings::Owned, error) {
    result := strings::Owned{}
    err := validate_path(path)
    if err == nil {
        err = preview_into(path, limit, &result.text)
    }
    return result, err
}

@[default_cc="c"]
extern {
    @[link_name="_telia_file_open"]
    fn file_open(path string, result *rawptr) error
    @[link_name="_telia_file_read"]
    fn file_read(state rawptr, buffer rawptr, size u64, count *u64, eof *bool) error
    @[link_name="_telia_file_close"]
    fn file_close(state *rawptr) error
    @[link_name="_telia_file_preview"]
    fn preview_into(path string, limit u64, result *string) error
}
