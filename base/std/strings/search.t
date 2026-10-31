package strings

fn contains(text string, needle string) bool {
    offset, found := find(text, needle)
    return found
}

fn starts_with(text string, prefix string) bool {
    if len(prefix) > len(text) {
        return false
    }
    candidate := slice(text, 0, len(prefix)) @fail
    return candidate == prefix
}

fn ends_with(text string, suffix string) bool {
    if len(suffix) > len(text) {
        return false
    }
    candidate := slice(text, len(text) - len(suffix), len(text)) @fail
    return candidate == suffix
}

fn find(text string, needle string) (u64, bool) {
    size := len(text)
    width := len(needle)
    if width > size {
        return 0, false
    }
    offset u64 := 0
    while offset <= size - width {
        candidate := slice(text, offset, offset + width) @fail
        if candidate == needle {
            return offset, true
        }
        if offset == size - width {
            return 0, false
        }
        offset = offset + 1
    }
    return 0, false
}
