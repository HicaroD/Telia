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
    InvalidPath
    FileOpenFailed
    FileReadFailed
    FileWriteFailed
    FileCloseFailed
    StreamReadFailed
    StreamWriteFailed
    StreamFlushFailed
    StreamModeFailed
    InvalidNumber
    InvalidBoolean
    InvalidEnvironmentName
    EnvironmentLookupFailed
}
