package process

use "std::unicode::utf8"
use "std::strings"
use "std::errors"

fn lookup_env(name string) (strings::Owned, bool, error) {
    result := strings::Owned{}
    found := false
    if strings::len(name) == 0 or strings::contains(name, "=") {
        return result, false, error(errors::Error::InvalidEnvironmentName, "invalid environment name")
    }
    if strings::contains(name, "\0") {
        return result, false, error(errors::Error::EmbeddedZero, "embedded zero in environment name")
    }
    validation_err := utf8::validate(name)
    if validation_err != nil {
        return result, false, validation_err
    }
    err := lookup_env_into(name, &result.text, &found)
    if err != nil {
        return result, false, err
    }
    value_err := utf8::validate(strings::view(&result))
    if value_err != nil {
        strings::free(&result)
        return result, false, value_err
    }
    return result, found, err
}

struct Args {
    state rawptr
}

fn args() (Args, error) {
    result := Args{}
    err := snapshot_into(&result.state)
    if err != nil {
        return result, err
    }
    index u64 := 0
    while index < count(&result) {
        text, access_err := get(&result, index)
        if access_err != nil {
            free(&result)
            return result, access_err
        }
        validation_err := utf8::validate(text)
        if validation_err != nil {
            free(&result)
            return result, validation_err
        }
        index = index + 1
    }
    return result, nil
}

fn count(snapshot *Args) u64 {
    return snapshot_count(snapshot.state)
}

fn get(snapshot *Args, index u64) (string, error) {
    result := ""
    err := get_into(snapshot.state, index, &result)
    return result, err
}

fn free(snapshot *Args) {
    free_state(&snapshot.state)
}

@[default_cc="c"]
extern {
    @[link_name="_telia_lookup_env"]
    fn lookup_env_into(name string, result *string, found *bool) error

    @[link_name="_telia_args_snapshot"]
    fn snapshot_into(result *rawptr) error
    @[link_name="_telia_args_count"]
    fn snapshot_count(state rawptr) u64
    @[link_name="_telia_args_get"]
    fn get_into(state rawptr, index u64, result *string) error
    @[link_name="_telia_args_free"]
    fn free_state(state *rawptr)
}
