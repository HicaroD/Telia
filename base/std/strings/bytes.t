package strings

use "std::lists"
use "std::slices"

// Byte access does not interpret UTF-8.
fn byte_at(text string, index u64) (u8, error) {
    value u32 := 0
    err := byte_into(text, index, &value)
    if err != nil {
        return 0, err
    }
    return narrow<u8>(value)
}

fn to_bytes(text string) (lists::List<u8>, error) {
    result := lists::create<u8>()
    index u64 := 0
    while index < len(text) {
        value, err := byte_at(text, index)
        if err == nil {
            err = lists::append<u8>(&result, value)
        }
        if err != nil {
            lists::free<u8>(&result)
            return result, err
        }
        index = index + 1
    }
    return result, nil
}

fn from_bytes(data u8[]) (Owned, error) {
    result := Owned{}
    err := bytes_clone(slice_data<u8>(data), slices::length<u8>(data), &result.text)
    return result, err
}

@[default_cc="c"]
extern {
    @[link_name="_telia_string_byte"]
    fn byte_into(text string, index u64, result *u32) error
    @[link_name="_telia_string_from_bytes"]
    fn bytes_clone(data rawptr, size u64, result *string) error
}
