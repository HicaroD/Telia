package maps

use "std::internal::swiss"
use "std::memory"

struct Map<K, V> {
    table swiss::Table<K, V>
}

struct Iterator<K, V> {
    source *Map<K, V>
    index u64
}

fn create<K, V>() Map<K, V> {
    return zero<Map<K, V>>()
}

fn create_in<K, V>(arena *memory::Arena) Map<K, V> {
    result := create<K, V>()
    result.table.arena = arena
    return result
}

fn clone<K, V>(map *Map<K, V>) (Map<K, V>, error) {
    return clone_in(map, nil)
}

fn clone_in<K, V>(map *Map<K, V>, arena *memory::Arena) (Map<K, V>, error) {
    table, err := swiss::clone_in(&map.table, arena)
    return Map<K, V>{table: table}, err
}

fn count<K, V>(map *Map<K, V>) u64 {
    return map.table.count
}

fn capacity<K, V>(map *Map<K, V>) u64 {
    return swiss::capacity(&map.table)
}

fn reserve<K, V>(map *Map<K, V>, requested u64) error {
    return swiss::reserve(&map.table, requested)
}

fn free<K, V>(map *Map<K, V>) {
    swiss::release(&map.table)
    return
}

fn clear<K, V>(map *Map<K, V>) {
    swiss::clear(&map.table)
    return
}

fn get_ptr<K, V>(map *Map<K, V>, key K) *V {
    index, found := swiss::find(&map.table, key)
    if not found {
        return nil
    }
    entries := map.table.entries
    return &entries[index].value
}

fn get<K, V>(map *Map<K, V>, key K) (V, bool) {
    pointer := get_ptr(map, key)
    if pointer == nil {
        return zero<V>(), false
    }
    return *pointer, true
}

fn contains<K, V>(map *Map<K, V>, key K) bool {
    return get_ptr(map, key) != nil
}

fn insert<K, V>(map *Map<K, V>, key K, value V) (V, bool, error) {
    pointer := get_ptr(map, key)
    if pointer != nil {
        old := *pointer
        *pointer = value
        return old, true, nil
    }
    err := swiss::insert_new(&map.table, key, value)
    if err != nil {
        return zero<V>(), false, err
    }
    return zero<V>(), false, nil
}

fn insert_if_absent<K, V>(map *Map<K, V>, key K, value V) (bool, error) {
    if contains(map, key) {
        return false, nil
    }
    err := swiss::insert_new(&map.table, key, value)
    if err != nil {
        return false, err
    }
    return true, nil
}

fn remove<K, V>(map *Map<K, V>, key K) (K, V, bool) {
    index, found := swiss::find(&map.table, key)
    if not found {
        return zero<K>(), zero<V>(), false
    }
    entries := map.table.entries
    controls := map.table.controls
    stored := entries[index]
    controls[index] = 129
    map.table.count = map.table.count - 1
    return stored.key, stored.value, true
}

fn iter<K, V>(map *Map<K, V>) Iterator<K, V> {
    return Iterator<K, V>{source: map, index: 0}
}

fn next<K, V>(iterator *Iterator<K, V>) (K, V, bool) {
    if iterator.source == nil {
        return zero<K>(), zero<V>(), false
    }
    controls := iterator.source.table.controls
    entries := iterator.source.table.entries
    while iterator.index < len(controls) {
        index := iterator.index
        iterator.index = index + 1
        if controls[index] < 128 {
            entry := entries[index]
            return entry.key, entry.value, true
        }
    }
    return zero<K>(), zero<V>(), false
}
