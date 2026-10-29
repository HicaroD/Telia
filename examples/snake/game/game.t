package game

const Width i32 = 28
const Height i32 = 24
const Capacity u64 = 672

enum Direction {
  Up
  Down
  Left
  Right
}

struct Cell {
  x i32
  y i32
}

struct State {
  body Cell[672]
  length u64
  direction Direction
  next_direction Direction
  turn_queued bool
  food Cell
  random u32
  score i32
  paused bool
  over bool
  won bool
}

fn same(left Cell, right Cell) bool {
  return left.x == right.x and left.y == right.y
}

fn occupied(state *State, cell Cell) bool {
  index u64 := 0
  while index < state.length {
    if same(state.body[index], cell) {
      return true
    }
    index = index + 1
  }
  return false
}

fn random_next(state *State) u32 {
  value := state.random
  value = value ^ (value << 13)
  value = value ^ (value >> 17)
  value = value ^ (value << 5)
  state.random = value
  return value
}

fn place_food(state *State) {
  if state.length == Capacity {
    state.won = true
    return
  }
  free_count u32 := 0
  y i32 := 0
  while y < Height {
    x i32 := 0
    while x < Width {
      if occupied(state, Cell{x: x, y: y}) == false {
        free_count = free_count + 1
      }
      x = x + 1
    }
    y = y + 1
  }
  value := random_next(state)
  target := value - (value / free_count) * free_count
  y = 0
  while y < Height {
    food_x i32 := 0
    while food_x < Width {
      cell := Cell{x: food_x, y: y}
      if occupied(state, cell) == false {
        if target == 0 {
          state.food = cell
          return
        }
        target = target - 1
      }
      food_x = food_x + 1
    }
    y = y + 1
  }
}

fn create(seed u32) State {
  state := State{length: 3, direction: Direction::Right, next_direction: Direction::Right, random: seed}
  if state.random == 0 {
    state.random = 1
  }
  state.body[0] = Cell{x: 14, y: 12}
  state.body[1] = Cell{x: 13, y: 12}
  state.body[2] = Cell{x: 12, y: 12}
  place_food(&state)
  return state
}

fn turn(state *State, direction Direction) {
  if state.over or state.won or state.paused or state.turn_queued {
    return
  }
  if direction == state.direction {
    return
  }
  if state.direction == Direction::Up and direction == Direction::Down {
    return
  }
  if state.direction == Direction::Down and direction == Direction::Up {
    return
  }
  if state.direction == Direction::Left and direction == Direction::Right {
    return
  }
  if state.direction == Direction::Right and direction == Direction::Left {
    return
  }
  state.next_direction = direction
  state.turn_queued = true
}

fn toggle_pause(state *State) {
  if state.over == false and state.won == false {
    state.paused = not state.paused
  }
}

fn step(state *State) {
  if state.over or state.won or state.paused {
    return
  }
  head := state.body[0]
  if state.next_direction == Direction::Up {
    head.y = head.y - 1
  } elif state.next_direction == Direction::Down {
    head.y = head.y + 1
  } elif state.next_direction == Direction::Left {
    head.x = head.x - 1
  } else {
    head.x = head.x + 1
  }
  if head.x < 0 or head.x >= Width or head.y < 0 or head.y >= Height {
    state.over = true
    return
  }
  growing := same(head, state.food)
  collision_length := state.length
  if growing == false {
    collision_length = collision_length - 1
  }
  index u64 := 0
  while index < collision_length {
    if same(head, state.body[index]) {
      state.over = true
      return
    }
    index = index + 1
  }
  if growing {
    state.length = state.length + 1
    state.score = state.score + 1
  }
  index = state.length - 1
  while index > 0 {
    state.body[index] = state.body[index - 1]
    index = index - 1
  }
  state.body[0] = head
  state.direction = state.next_direction
  state.turn_queued = false
  if growing {
    place_food(state)
  }
}
