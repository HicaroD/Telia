package swiss

use "std::memory"
use "std::errors"

// Internal shared engine. Controls: 0..127 live fingerprint, 128 empty,
// 129 deleted. Eight-bucket groups, triangular group probing, 7/8 load.
struct Entry<K, V> {
    key K
    value V
}

struct Table<K, V> {
    controls u8[]
    entries Entry<K, V>[]
    count u64
    occupied u64
    seed u64[2]
    arena *memory::Arena
}

fn capacity<K, V>(table *Table<K, V>) u64 {
    return len(table.controls) / 8 * 7
}

fn hash<K, V>(table *Table<K, V>, key K) u64 {
    return __key_hash<K>(key, table.seed[0], table.seed[1])
}

fn fingerprint(hash u64) u8 {
    result := narrow<u8>(hash - hash / 128 * 128) @fail
    return result
}

fn first_group(hash u64, groups u64) u64 {
    high := hash / 128
    return high - high / groups * groups
}

fn advance(group u64, step u64, groups u64) u64 {
    sum := group + step
    return sum - sum / groups * groups
}

fn find<K, V>(table *Table<K, V>, key K) (u64, bool) {
    if table.count == 0 {
        return 0, false
    }
    hashed := hash(table, key)
    tag := fingerprint(hashed)
    controls := table.controls
    entries := table.entries
    groups := len(controls) / 8
    group := first_group(hashed, groups)
    step u64 := 0
    while step < groups {
        empty := false
        lane u64 := 0
        while lane < 8 {
            index := group * 8 + lane
            if controls[index] == tag {
                if entries[index].key == key {
                    return index, true
                }
            }
            if controls[index] == 128 {
                empty = true
            }
            lane = lane + 1
        }
        if empty {
            return 0, false
        }
        step = step + 1
        group = advance(group, step, groups)
    }
    return 0, false
}

// No equality here: rebuilding must preserve even non-reflexive NaN keys.
// Caller guarantees at least one empty/deleted slot and a committed seed.
fn vacant<K, V>(table *Table<K, V>, hashed u64) u64 {
    controls := table.controls
    groups := len(controls) / 8
    group := first_group(hashed, groups)
    step u64 := 0
    while step < groups {
        lane u64 := 0
        while lane < 8 {
            index := group * 8 + lane
            if controls[index] >= 128 {
                return index
            }
            lane = lane + 1
        }
        step = step + 1
        group = advance(group, step, groups)
    }
    // Unreachable under the 7/8 occupied bound.
    return len(controls)
}

fn release<K, V>(table *Table<K, V>) {
    if table.arena == nil {
        memory::free(slice_data<u8>(table.controls))
        memory::free(slice_data<Entry<K, V>>(table.entries))
    }
    arena := table.arena
    *table = zero<Table<K, V>>()
    table.arena = arena
    return
}

fn allocate<T>(arena *memory::Arena, buckets u64) (T[], error) {
    if buckets > 18446744073709551615 / sizeof<T>() {
        return zero<T[]>(), error(errors::Error::SizeOverflow, "table byte size overflow")
    }
    storage rawptr := nil
    if arena == nil {
        data, err := memory::alloc(buckets * sizeof<T>())
        if err != nil {
            return zero<T[]>(), err
        }
        storage = data
    } else {
        data, err := memory::arena_alloc(arena, buckets * sizeof<T>(), alignof<T>())
        if err != nil {
            return zero<T[]>(), err
        }
        storage = data
    }
    return slice_from_raw<T>(storage, buckets), nil
}

// Transactional rebuild: validate all arithmetic, seed, allocate both buffers,
// then copy live entries without equality. Publish only after all fallible work.
fn rebuild<K, V>(table *Table<K, V>, buckets u64) error {
    if buckets > 18446744073709551615 / sizeof<Entry<K, V>>() {
        return error(errors::Error::SizeOverflow, "table byte size overflow")
    }
    seed := table.seed
    if len(table.controls) == 0 {
        acquired, err := __key_seed()
        if err != nil {
            return err
        }
        seed = acquired
    }
    controls, control_error := allocate<u8>(table.arena, buckets)
    if control_error != nil {
        return control_error
    }
    entries, entry_error := allocate<Entry<K, V>>(table.arena, buckets)
    if entry_error != nil {
        if table.arena == nil {
            memory::free(slice_data<u8>(controls))
        }
        return entry_error
    }
    index u64 := 0
    while index < buckets {
        controls[index] = 128
        index = index + 1
    }
    fresh := zero<Table<K, V>>()
    fresh.controls = controls
    fresh.entries = entries
    fresh.seed = seed
    fresh.arena = table.arena
    old_controls := table.controls
    old_entries := table.entries
    index = 0
    while index < len(old_controls) {
        if old_controls[index] < 128 {
            entry := old_entries[index]
            hashed := hash(&fresh, entry.key)
            target := vacant(&fresh, hashed)
            entries[target] = entry
            controls[target] = fingerprint(hashed)
            fresh.count = fresh.count + 1
        }
        index = index + 1
    }
    fresh.occupied = fresh.count
    release(table)
    *table = fresh
    return nil
}

fn reserve<K, V>(table *Table<K, V>, requested u64) error {
    if requested <= capacity(table) {
        return nil
    }
    buckets u64 := 8
    while buckets / 8 * 7 < requested {
        // Also bounds group-probe addition independently of element layout.
        if buckets > 9223372036854775807 {
            return error(errors::Error::SizeOverflow, "table capacity overflow")
        }
        buckets = buckets * 2
    }
    return rebuild(table, buckets)
}

fn prepare<K, V>(table *Table<K, V>) error {
    if table.count == 18446744073709551615 {
        return error(errors::Error::SizeOverflow, "table count overflow")
    }
    if table.count + 1 > capacity(table) {
        return reserve(table, table.count + 1)
    }
    if table.occupied + 1 > capacity(table) {
        return rebuild(table, len(table.controls))
    }
    return nil
}

fn insert_new<K, V>(table *Table<K, V>, key K, value V) error {
    prepare(table) @try
    hashed := hash(table, key)
    index := vacant(table, hashed)
    controls := table.controls
    entries := table.entries
    if controls[index] == 128 {
        table.occupied = table.occupied + 1
    }
    entries[index] = Entry<K, V>{key: key, value: value}
    controls[index] = fingerprint(hashed)
    table.count = table.count + 1
    return nil
}

fn clear<K, V>(table *Table<K, V>) {
    controls := table.controls
    index u64 := 0
    while index < len(controls) {
        controls[index] = 128
        index = index + 1
    }
    table.count = 0
    table.occupied = 0
    return
}
