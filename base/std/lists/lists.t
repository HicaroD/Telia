package lists

use "std::memory"

struct List<T> {
    buffer T[]
    count u64
    arena *memory::Arena
}

fn create<T>() List<T> {
    return zero<List<T>>()
}

fn create_in<T>(arena *memory::Arena) List<T> {
    result := zero<List<T>>()
    result.arena = arena
    return result
}

fn count<T>(list *List<T>) u64 {
    return list.count
}

fn capacity<T>(list *List<T>) u64 {
    return len(list.buffer)
}

fn view<T>(list *List<T>) T[] {
    buffer := list.buffer
    result := buffer[0:list.count] @fail
    return result
}

fn free<T>(list *List<T>) {
    if list.arena == nil {
        memory::free(slice_data<T>(list.buffer))
    }
    list.buffer = zero<T[]>()
    list.count = 0
    return
}

fn allocate_buffer<T>(arena *memory::Arena, requested u64) (T[], error) {
    if requested > 18446744073709551615 / sizeof<T>() {
        return zero<T[]>(), error("list capacity overflow")
    }
    storage rawptr := nil
    if arena == nil {
        data, err := memory::alloc(requested * sizeof<T>())
        if err != nil {
            return zero<T[]>(), err
        }
        storage = data
    } else {
        data, err := memory::arena_alloc(arena, requested * sizeof<T>(), alignof<T>())
        if err != nil {
            return zero<T[]>(), err
        }
        storage = data
    }
    buffer := slice_from_raw<T>(storage, requested)
    return buffer, nil
}

fn reserve<T>(list *List<T>, requested u64) error {
    if requested <= len(list.buffer) {
        return nil
    }
    buffer, err := allocate_buffer<T>(list.arena, requested)
    if err != nil {
        return err
    }
    old := list.buffer
    index u64 := 0
    while index < list.count {
        buffer[index] = old[index]
        index = index + 1
    }
    if list.arena == nil {
        memory::free(slice_data<T>(old))
    }
    list.buffer = buffer
    return nil
}

fn append<T>(list *List<T>, value T) error {
    single T[1] := .{value}
    return append_slice<T>(list, single[:])
}

fn append_slice<T>(list *List<T>, values T[]) error {
    if len(values) > 18446744073709551615 - list.count {
        return error("list length overflow")
    }
    needed := list.count + len(values)
    buffer := list.buffer
    if needed > len(list.buffer) {
        target := needed
        current := len(list.buffer)
        if current <= (18446744073709551615 / sizeof<T>()) / 2 {
            doubled := current * 2
            if doubled > target {
                target = doubled
            }
        }
        grown, err := allocate_buffer<T>(list.arena, target)
        if err != nil {
            return err
        }
        index u64 := 0
        while index < list.count {
            grown[index] = buffer[index]
            index = index + 1
        }
        buffer = grown
    }
    index u64 := 0
    while index < len(values) {
        buffer[list.count + index] = values[index]
        index = index + 1
    }
    if needed > len(list.buffer) {
        if list.arena == nil {
            memory::free(slice_data<T>(list.buffer))
        }
        list.buffer = buffer
    }
    list.count = needed
    return nil
}

fn get<T>(list *List<T>, index u64) (T, error) {
    if index >= list.count {
        return zero<T>(), error("list index out of bounds")
    }
    buffer := list.buffer
    return buffer[index], nil
}

fn set<T>(list *List<T>, index u64, value T) error {
    if index >= list.count {
        return error("list index out of bounds")
    }
    buffer := list.buffer
    buffer[index] = value
    return nil
}

fn pop<T>(list *List<T>) (T, bool) {
    if list.count == 0 {
        return zero<T>(), false
    }
    list.count = list.count - 1
    buffer := list.buffer
    return buffer[list.count], true
}

fn remove<T>(list *List<T>, index u64) (T, error) {
    if index >= list.count {
        return zero<T>(), error("list index out of bounds")
    }
    buffer := list.buffer
    result := buffer[index]
    cursor := index
    while cursor + 1 < list.count {
        buffer[cursor] = buffer[cursor + 1]
        cursor = cursor + 1
    }
    list.count = list.count - 1
    return result, nil
}

fn remove_unordered<T>(list *List<T>, index u64) (T, error) {
    if index >= list.count {
        return zero<T>(), error("list index out of bounds")
    }
    buffer := list.buffer
    result := buffer[index]
    list.count = list.count - 1
    buffer[index] = buffer[list.count]
    return result, nil
}

fn clear<T>(list *List<T>) {
    list.count = 0
    return
}

fn clone<T>(list *List<T>) (List<T>, error) {
    return clone_in<T>(list, nil)
}

fn clone_in<T>(list *List<T>, arena *memory::Arena) (List<T>, error) {
    result := create_in<T>(arena)
    source := view<T>(list)
    err := append_slice<T>(&result, source)
    return result, err
}
