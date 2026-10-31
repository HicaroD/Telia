package path

use "std::strings"
use "std::errors"
use "std::unicode::utf8"

fn validate(path string) error {
    if strings::len(path) == 0 {
        return error(errors::Error::InvalidPath, "empty path")
    }
    if strings::contains(path, "\0") {
        return error(errors::Error::EmbeddedZero, "embedded zero in path")
    }
    return utf8::validate(path)
}

fn canonicalize(path string) (strings::Owned, error) {
    return resolve(path, false)
}

// Only the final component may be missing; its parent must exist.
fn canonicalize_for_create(path string) (strings::Owned, error) {
    return resolve(path, true)
}

fn resolve(path string, allow_missing bool) (strings::Owned, error) {
    result := strings::Owned{}
    err := validate(path)
    if err == nil {
        err = canonical_into(path, allow_missing, &result.text)
    }
    if err == nil {
        err = utf8::validate(strings::view(&result))
    }
    if err != nil {
        strings::free(&result)
    }
    return result, err
}

@[default_cc="c"]
extern {
    @[link_name="_telia_path_canonical"]
    fn canonical_into(path string, allow_missing bool, result *string) error
}
