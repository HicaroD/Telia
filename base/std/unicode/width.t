package unicode

use "std::errors"

// Code-point cell width, not grapheme/emoji-sequence or tab layout.
fn width(point u32) (i32, error) {
    if point > 1114111 or (point >= 55296 and point <= 57343) {
        return 0, error(errors::Error::InvalidCodePoint, "invalid Unicode scalar value")
    }
    if point == 0 {
        return 0, nil
    }
    if nonprintable(point) {
        return -1, nil
    }
    if zero_width(point) {
        return 0, nil
    }
    if wide(point) {
        return 2, nil
    }
    return 1, nil
}

fn in_ranges(point u32, ranges u32[]) bool {
    low u64 := 0
    high := len(ranges) / 2
    while low < high {
        middle := low + (high - low) / 2
        if point < ranges[middle * 2] {
            high = middle
        } elif point > ranges[middle * 2 + 1] {
            low = middle + 1
        } else {
            return true
        }
    }
    return false
}
