package sets

use "std::internal::swiss"
use "std::memory"

// The C backend gives empty records a nonzero representation. This concrete
// value is internal: sets share the table engine, not a second probing algorithm.
struct Unit {}

struct Set<T> {
    table swiss::Table<T, Unit>
}

struct Iterator<T> {
    source *Set<T>
    index u64
}

fn create<T>() Set<T> {
    return zero<Set<T>>()
}

fn create_in<T>(arena *memory::Arena) Set<T> {
    result := create<T>()
    result.table.arena = arena
    return result
}

fn clone<T>(set *Set<T>) (Set<T>, error) {
    return clone_in(set, nil)
}

fn clone_in<T>(set *Set<T>, arena *memory::Arena) (Set<T>, error) {
    table, err := swiss::clone_in(&set.table, arena)
    return Set<T>{table: table}, err
}

fn count<T>(set *Set<T>) u64 {
    return set.table.count
}

fn capacity<T>(set *Set<T>) u64 {
    return swiss::capacity(&set.table)
}

fn reserve<T>(set *Set<T>, requested u64) error {
    return swiss::reserve(&set.table, requested)
}

fn free<T>(set *Set<T>) {
    swiss::release(&set.table)
    return
}

fn clear<T>(set *Set<T>) {
    swiss::clear(&set.table)
    return
}

fn contains<T>(set *Set<T>, element T) bool {
    index, found := swiss::find(&set.table, element)
    return found
}

fn insert<T>(set *Set<T>, element T) (bool, error) {
    index, found := swiss::find(&set.table, element)
    if found {
        return false, nil
    }
    err := swiss::insert_new(&set.table, element, Unit{})
    if err != nil {
        return false, err
    }
    return true, nil
}

fn remove<T>(set *Set<T>, element T) (T, bool) {
    index, found := swiss::find(&set.table, element)
    if not found {
        return zero<T>(), false
    }
    entries := set.table.entries
    controls := set.table.controls
    stored := entries[index].key
    controls[index] = 129
    set.table.count = set.table.count - 1
    return stored, true
}

fn iter<T>(set *Set<T>) Iterator<T> {
    return Iterator<T>{source: set, index: 0}
}

fn next<T>(iterator *Iterator<T>) (T, bool) {
    if iterator.source == nil {
        return zero<T>(), false
    }
    controls := iterator.source.table.controls
    entries := iterator.source.table.entries
    while iterator.index < len(controls) {
        index := iterator.index
        iterator.index = index + 1
        if controls[index] < 128 {
            return entries[index].key, true
        }
    }
    return zero<T>(), false
}
