package slices

use "std::errors"

fn length<T>(values T[]) u64 {
    return len(values)
}

fn get<T>(values T[], index u64) (T, error) {
    if index >= len(values) {
        return zero<T>(), error(errors::Error::OutOfBounds, "slice index out of bounds")
    }
    return values[index], nil
}

fn set<T>(values T[], index u64, value T) error {
    if index >= len(values) {
        return error(errors::Error::OutOfBounds, "slice index out of bounds")
    }
    values[index] = value
    return nil
}
