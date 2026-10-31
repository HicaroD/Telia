package utf8

use "std::strings"

struct Iterator {
    text string
    offset u64
}

fn iterator(text string) (Iterator, error) {
    err := validate(text)
    if err != nil {
        return Iterator{}, err
    }
    return Iterator{text: text, offset: 0}, nil
}

fn next(iter *Iterator) (u32, u64, bool, error) {
    if iter.offset == strings::len(iter.text) {
        return 0, 0, false, nil
    }
    point, width, err := decode(iter.text, iter.offset)
    if err != nil {
        return 0, 0, false, err
    }
    offset := iter.offset
    iter.offset = offset + width
    return point, offset, true, nil
}

@[default_cc="c"]
extern {
    @[link_name="_telia_string_byte"]
    fn byte_into(text string, index u64, result *u32) error
}

fn read_byte(text string, index u64) (u32, error) {
    result u32 := 0
    err := byte_into(text, index, &result)
    return result, err
}

fn decode(text string, offset u64) (u32, u64, error) {
    size := strings::len(text)
    if offset >= size {
        return 0, 0, error("invalid UTF-8 decode offset")
    }
    first, first_err := read_byte(text, offset)
    if first_err != nil {
        return 0, 0, first_err
    }
    if first < 128 {
        return first, 1, nil
    }
    width u64 := 0
    point u32 := 0
    minimum u32 := 0
    if first >= 194 and first <= 223 {
        width = 2
        point = first & 31
        minimum = 128
    } elif first >= 224 and first <= 239 {
        width = 3
        point = first & 15
        minimum = 2048
    } elif first >= 240 and first <= 244 {
        width = 4
        point = first & 7
        minimum = 65536
    } else {
        return 0, 0, error("invalid UTF-8")
    }
    if width > size - offset {
        return 0, 0, error("truncated UTF-8")
    }
    index u64 := 1
    while index < width {
        next, next_err := read_byte(text, offset + index)
        if next_err != nil {
            return 0, 0, next_err
        }
        if next < 128 or next > 191 {
            return 0, 0, error("invalid UTF-8 continuation")
        }
        point = (point << 6) | (next & 63)
        index = index + 1
    }
    if point < minimum or point > 1114111 or (point >= 55296 and point <= 57343) {
        return 0, 0, error("invalid UTF-8 code point")
    }
    return point, width, nil
}

fn validate(text string) error {
    total, err := count(text)
    return err
}

fn count(text string) (u64, error) {
    total u64 := 0
    offset u64 := 0
    while offset < strings::len(text) {
        point, width, err := decode(text, offset)
        if err != nil {
            return 0, err
        }
        offset = offset + width
        total = total + 1
    }
    return total, nil
}
