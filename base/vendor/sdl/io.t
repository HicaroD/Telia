package sdl

use "vendor::sdl::raw"

fn io_from_bytes(bytes u8[]) (*IOStream, error) {
  if len(bytes) == 0 {
    return nil, error("SDL memory IO requires nonempty storage")
  }
  stream := raw::io_from_bytes_pointer(&bytes[0], len(bytes))
  if stream == nil {
    return nil, error("SDL io_from_bytes failed")
  }
  return stream, nil
}

fn read_bytes(stream *IOStream, bytes u8[]) u64 {
  if len(bytes) == 0 {
    return 0
  }
  return raw::read_bytes_pointer(stream, &bytes[0], len(bytes))
}
