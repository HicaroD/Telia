package memory

struct Arena {
  state rawptr
}

@[default_cc="c"]
extern {
  @[link_name="_telia_memory_alloc"]
  fn alloc_into(size u64, result *rawptr) error

  @[link_name="_telia_memory_free"]
  fn free(ptr rawptr)

  @[link_name="_telia_arena_alloc"]
  fn arena_alloc_into(state *rawptr, size u64, alignment u64, result *rawptr) error

  @[link_name="_telia_arena_reset"]
  fn reset_state(state rawptr)

  @[link_name="_telia_arena_destroy"]
  fn destroy_state(state *rawptr)
}

fn alloc(size u64) (rawptr, error) {
  result rawptr := nil
  err := alloc_into(size, &result)
  return result, err
}

fn arena_create() (Arena, error) {
  result := Arena.{state: nil}
  return result, nil
}

fn arena_alloc(arena *Arena, size u64, alignment u64) (rawptr, error) {
  result rawptr := nil
  err := arena_alloc_into(&arena.state, size, alignment, &result)
  return result, err
}

fn arena_reset(arena *Arena) {
  reset_state(arena.state)
  return
}

fn arena_destroy(arena *Arena) {
  destroy_state(&arena.state)
  return
}
