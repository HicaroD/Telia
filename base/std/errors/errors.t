package errors

enum Error: error {
    OutOfBounds
    InvalidRange
    SizeOverflow
    AllocationFailed
    InvalidAlignment
    EmbeddedZero
    EmptySeparator
    InvalidUTF8
    ConversionOutOfRange
}
