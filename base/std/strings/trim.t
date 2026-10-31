package strings

use "std::unicode::utf8"

fn trim_space(text string) (string, error) {
    start := len(text)
    end u64 := 0
    offset u64 := 0
    while offset < len(text) {
        point, width, err := utf8::decode(text, offset)
        if err != nil {
            return "", err
        }
        if not is_white_space(point) {
            if start == len(text) {
                start = offset
            }
            end = offset + width
        }
        offset = offset + width
    }
    if end == 0 {
        return "", nil
    }
    return slice(text, start, end)
}

fn is_white_space(point u32) bool {
    if (point >= 9 and point <= 13) or point == 32 {
        return true
    }
    if point == 133 or point == 160 or point == 5760 {
        return true
    }
    if point >= 8192 and point <= 8202 {
        return true
    }
    return point == 8232 or point == 8233 or point == 8239 or point == 8287 or point == 12288
}
