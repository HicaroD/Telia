package slices

fn get<T>(values T[], index u64) (T, error) {
    if index >= len(values) {
        return zero<T>(), error("slice index out of bounds")
    }
    return values[index], nil
}

fn set<T>(values T[], index u64, value T) error {
    if index >= len(values) {
        return error("slice index out of bounds")
    }
    values[index] = value
    return nil
}
