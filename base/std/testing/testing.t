package testing

@[default_cc="c"]
extern {
    @[link_name="_telia_testing_fail"]
    fn fail()
}

fn assert(condition bool) {
    if not condition {
        fail()
    }
}

fn assert_eq<T>(actual T, expected T) {
    assert(actual == expected)
}

fn assert_ne<T>(actual T, unexpected T) {
    assert(actual != unexpected)
}
