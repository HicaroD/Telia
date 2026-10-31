package fs

use "std::strings"
use "std::errors"
use "std::unicode::utf8"

fn validate_path(path string) error {
    if strings::len(path) == 0 {
        return error(errors::Error::InvalidPath, "empty file path")
    }
    if strings::contains(path, "\0") {
        return error(errors::Error::EmbeddedZero, "embedded zero in file path")
    }
    return utf8::validate(path)
}

fn read_file(path string) (strings::Owned, error) {
    result := strings::Owned{}
    path_err := validate_path(path)
    if path_err != nil {
        return result, path_err
    }
    err := read_into(path, &result.text)
    return result, err
}

fn write_file(path string, data string) error {
    path_err := validate_path(path)
    if path_err != nil {
        return path_err
    }
    return write_bytes(path, data)
}

@[default_cc="c"]
extern {
    @[link_name="_telia_fs_read"]
    fn read_into(path string, result *string) error

    @[link_name="_telia_fs_write"]
    fn write_bytes(path string, data string) error
}
