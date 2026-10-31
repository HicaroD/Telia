package strconv

use "std::strings"
use "std::errors"

fn format_i64(value i64) (strings::Owned, error) {
    negative := value < 0
    magnitude u64 := 0
    if negative {
        // Avoid negating the signed minimum itself.
        magnitude = narrow<u64>(-(value + 1)) @try
        magnitude = magnitude + 1
    } else {
        magnitude = narrow<u64>(value) @try
    }
    result := strings::Owned{text: ""}
    err := format_decimal_into(magnitude, negative, &result.text)
    return result, err
}

fn format_bool(value bool) (strings::Owned, error) {
    if value {
        return strings::clone("true")
    }
    return strings::clone("false")
}

fn format_u64(value u64) (strings::Owned, error) {
    result := strings::Owned{text: ""}
    err := format_decimal_into(value, false, &result.text)
    return result, err
}

@[default_cc="c"]
extern {
    @[link_name="_telia_format_decimal"]
    fn format_decimal_into(magnitude u64, negative bool, result *string) error

    @[link_name="_telia_string_byte"]
    fn byte_into(text string, index u64, result *u32) error
}

fn parse_u64(text string) (u64, error) {
    return parse_magnitude(text, 0, 18446744073709551615)
}

fn parse_bool(text string) (bool, error) {
    if text == "true" {
        return true, nil
    }
    if text == "false" {
        return false, nil
    }
    return false, error(errors::Error::InvalidBoolean, "expected true or false")
}

fn parse_i64(text string) (i64, error) {
    start u64 := 0
    negative := false
    if strings::len(text) != 0 {
        first u32 := 0
        err := byte_into(text, 0, &first)
        if err != nil {
            return 0, err
        }
        negative = first == 45
        if negative or first == 43 {
            start = 1
        }
    }
    limit u64 := 9223372036854775807
    if negative {
        limit = limit + 1
    }
    magnitude, err := parse_magnitude(text, start, limit)
    if err != nil {
        return 0, err
    }
    if negative {
        if magnitude == 9223372036854775808 {
            return -9223372036854775807 - 1, nil
        }
        value := narrow<i64>(magnitude) @try
        return -value, nil
    }
    return narrow<i64>(magnitude)
}

fn parse_magnitude(text string, start u64, limit u64) (u64, error) {
    size := strings::len(text)
    if start == size {
        return 0, error(errors::Error::InvalidNumber, "expected decimal digits")
    }
    value u64 := 0
    index := start
    while index < size {
        byte u32 := 0
        err := byte_into(text, index, &byte)
        if err != nil {
            return 0, err
        }
        if byte < 48 or byte > 57 {
            return 0, error(errors::Error::InvalidNumber, "expected decimal digits")
        }
        digit := widen<u64>(byte - 48)
        if value > (limit - digit) / 10 {
            return 0, error(errors::Error::ConversionOutOfRange, "integer out of range")
        }
        value = value * 10 + digit
        index = index + 1
    }
    return value, nil
}
