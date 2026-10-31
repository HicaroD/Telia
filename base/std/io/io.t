package io

use "std::strings"

fn read_stdin() (strings::Owned, error) {
    result := strings::Owned{}
    err := read_into(&result.text)
    return result, err
}

fn write_stdout(data string) error {
    return write_bytes(data, false)
}

fn write_stderr(data string) error {
    return write_bytes(data, true)
}

@[default_cc="c"]
extern {
  @[link_name="_telia_io_read"]
  fn read_into(result *string) error

  @[link_name="_telia_io_write"]
  fn write_bytes(data string, stderr bool) error

  @[link_name="_telia_println"]
  fn println(text string)
}
