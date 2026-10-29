package strings

use "std::memory"
use "std::ffi::c"

struct Owned {
  text string
}

struct COwned {
  text string
}

@[default_cc="c"]
extern {
  @[link_name="_telia_string_len"]
  fn len(text string) u64

  @[link_name="_telia_string_slice"]
  fn slice_into(text string, start u64, end u64, result *string) error

  @[link_name="_telia_string_clone"]
  fn clone_into(text string, result *string) error

  @[link_name="_telia_string_concat"]
  fn concat_into(left string, right string, result *string) error

  @[link_name="_telia_string_free"]
  fn free_text(text *string)

  @[link_name="_telia_string_to_c"]
  fn to_c_into(text string, result *string) error

  @[link_name="_telia_string_c_view"]
  fn c_pointer(text string) c::CString

  @[link_name="_telia_string_from_c"]
  fn from_c(text c::ConstCString) string

  @[link_name="_telia_arena_string_concat"]
  fn concat_in_into(state *rawptr, left string, right string, result *string) error

  @[link_name="_telia_arena_string_to_c"]
  fn to_c_in_into(state *rawptr, text string, result *string) error
}

fn clone(text string) (Owned, error) {
  result := Owned{text: ""}
  err := clone_into(text, &result.text)
  return result, err
}

fn concat(left string, right string) (Owned, error) {
  result := Owned{text: ""}
  err := concat_into(left, right, &result.text)
  return result, err
}

fn view(owner *Owned) string {
  return owner.text
}

fn free(owner *Owned) {
  free_text(&owner.text)
  return
}

fn to_c(text string) (COwned, error) {
  result := COwned{text: ""}
  err := to_c_into(text, &result.text)
  return result, err
}

fn c_view(owner *COwned) c::CString {
  return c_pointer(owner.text)
}

fn free_c(owner *COwned) {
  free_text(&owner.text)
  return
}

fn clone_in(arena *memory::Arena, text string) (string, error) {
  result := ""
  err := concat_in_into(&arena.state, text, "", &result)
  return result, err
}

fn concat_in(arena *memory::Arena, left string, right string) (string, error) {
  result := ""
  err := concat_in_into(&arena.state, left, right, &result)
  return result, err
}

fn to_c_in(arena *memory::Arena, text string) (c::CString, error) {
  result := ""
  err := to_c_in_into(&arena.state, text, &result)
  return c_pointer(result), err
}

fn slice(text string, start u64, end u64) (string, error) {
  result := ""
  err := slice_into(text, start, end, &result)
  return result, err
}
