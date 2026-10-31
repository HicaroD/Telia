package strings

use "std::errors"

use "std::lists"

fn split(text string, separator string) (lists::List<string>, error) {
    result := lists::create<string>()
    width := len(separator)
    if width == 0 {
        return result, error(errors::Error::EmptySeparator, "empty split separator")
    }
    start u64 := 0
    while true {
        remaining := slice(text, start, len(text)) @fail
        offset, found := find(remaining, separator)
        end := len(text)
        if found {
            end = start + offset
        }
        field := slice(text, start, end) @fail
        err := lists::append<string>(&result, field)
        if err != nil {
            lists::free<string>(&result)
            return result, err
        }
        if not found {
            return result, nil
        }
        start = end + width
    }
    return result, nil
}
