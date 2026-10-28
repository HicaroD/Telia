// Private implementation of the standard library's raw-memory operations.
// Public contracts live in base/std; this file uses only the primitive ABI.
#define _TELIA_ERROR(message) ((_Error){{(const unsigned char *)(message), sizeof(message) - 1}})

// Override only at the native build boundary for deterministic failure tests.
#ifndef _TELIA_MALLOC
#define _TELIA_MALLOC malloc
#endif
#ifndef _TELIA_FREE
#define _TELIA_FREE free
#endif

static inline uint64_t _telia_string_len(_String text) {
    return text.len;
}

static inline _Error _telia_string_slice(_String text, uint64_t start,
                                        uint64_t end, _String *result) {
    *result = (_String){0};
    if (start > end || end > text.len) {
        return _TELIA_ERROR("invalid string range");
    }
    if (start != end) {
        *result = (_String){text.data + (size_t)start, (size_t)(end - start)};
    }
    return (_Error){0};
}

static inline _Error _telia_string_concat(_String left, _String right, _String *result) {
    *result = (_String){0};
    if (right.len > SIZE_MAX - left.len) return _TELIA_ERROR("allocation size overflow");
    size_t size = left.len + right.len;
    if (!size) return (_Error){0};
    unsigned char *data = _TELIA_MALLOC(size);
    if (!data) return _TELIA_ERROR("allocation failed");
    if (left.len) memcpy(data, left.data, left.len);
    if (right.len) memcpy(data + left.len, right.data, right.len);
    *result = (_String){data, size};
    return (_Error){0};
}

static inline _Error _telia_string_clone(_String text, _String *result) {
    return _telia_string_concat(text, (_String){0}, result);
}

static inline void _telia_string_free(_String *text) {
    _TELIA_FREE((void *)text->data);
    *text = (_String){0};
}

static inline _Error _telia_string_to_c(_String text, _String *result) {
    *result = (_String){0};
    if (text.len == SIZE_MAX) return _TELIA_ERROR("allocation size overflow");
    if (text.len && memchr(text.data, 0, text.len)) return _TELIA_ERROR("embedded zero byte");
    unsigned char *data = _TELIA_MALLOC(text.len + 1);
    if (!data) return _TELIA_ERROR("allocation failed");
    if (text.len) memcpy(data, text.data, text.len);
    data[text.len] = 0;
    *result = (_String){data, text.len};
    return (_Error){0};
}

static inline char *_telia_string_c_view(_String text) {
    return (char *)text.data;
}

static inline _String _telia_string_from_c(const char *text) {
    return (_String){(const unsigned char *)text, strlen(text)};
}

static inline _Error _telia_memory_alloc(uint64_t size, void **result) {
    *result = NULL;
    if (size > SIZE_MAX) return _TELIA_ERROR("allocation size overflow");
    if (!size) return (_Error){0};
    void *data = _TELIA_MALLOC((size_t)size);
    if (!data) return _TELIA_ERROR("allocation failed");
    *result = data;
    return (_Error){0};
}

static inline void _telia_memory_free(void *ptr) {
    _TELIA_FREE(ptr);
}

typedef struct _TeliaArenaBlock {
    struct _TeliaArenaBlock *next;
    size_t capacity;
    size_t used;
    unsigned char data[];
} _TeliaArenaBlock;

static inline _Error _telia_arena_alloc(void **state, uint64_t size,
                                       uint64_t alignment, void **result) {
    *result = NULL;
    if (!alignment || alignment > SIZE_MAX || (alignment & (alignment - 1))) {
        return _TELIA_ERROR("invalid allocation alignment");
    }
    if (size > SIZE_MAX || size > SIZE_MAX - (alignment - 1)) {
        return _TELIA_ERROR("allocation size overflow");
    }
    if (!size) return (_Error){0};
    _TeliaArenaBlock *last = NULL;
    for (_TeliaArenaBlock *block = *state; block; block = block->next) {
        size_t padding = (-(uintptr_t)(block->data + block->used)) & (size_t)(alignment - 1);
        size_t available = block->capacity - block->used;
        if (padding <= available && size <= available - padding) {
            unsigned char *data = block->data + block->used + padding;
            block->used += padding + (size_t)size;
            *result = data;
            return (_Error){0};
        }
        last = block;
    }
    size_t capacity = (size_t)size + ((size_t)alignment - 1);
    if (capacity < 4096) capacity = 4096;
    if (capacity > SIZE_MAX - sizeof(_TeliaArenaBlock)) {
        return _TELIA_ERROR("allocation size overflow");
    }
    _TeliaArenaBlock *block = _TELIA_MALLOC(sizeof(*block) + capacity);
    if (!block) return _TELIA_ERROR("allocation failed");
    block->next = NULL;
    block->capacity = capacity;
    size_t padding = (-(uintptr_t)block->data) & (size_t)(alignment - 1);
    block->used = padding + (size_t)size;
    *result = block->data + padding;
    // Publish only after allocation and initialization succeed.
    if (last) last->next = block;
    else *state = block;
    return (_Error){0};
}

static inline void _telia_arena_reset(void *state) {
    for (_TeliaArenaBlock *block = state; block; block = block->next) block->used = 0;
}

static inline void _telia_arena_destroy(void **state) {
    _TeliaArenaBlock *block = *state;
    while (block) {
        _TeliaArenaBlock *next = block->next;
        _TELIA_FREE(block);
        block = next;
    }
    *state = NULL;
}

static inline _Error _telia_arena_string_concat(void **state, _String left,
                                               _String right, _String *result) {
    *result = (_String){0};
    if (right.len > SIZE_MAX - left.len) return _TELIA_ERROR("allocation size overflow");
    size_t size = left.len + right.len;
    void *storage;
    _Error err = _telia_arena_alloc(state, size, 1, &storage);
    if (err.msg.data) return err;
    unsigned char *data = storage;
    if (left.len) memcpy(data, left.data, left.len);
    if (right.len) memcpy(data + left.len, right.data, right.len);
    *result = (_String){data, size};
    return (_Error){0};
}

static inline _Error _telia_arena_string_to_c(void **state, _String text, _String *result) {
    *result = (_String){0};
    if (text.len == SIZE_MAX) return _TELIA_ERROR("allocation size overflow");
    if (text.len && memchr(text.data, 0, text.len)) return _TELIA_ERROR("embedded zero byte");
    void *storage;
    _Error err = _telia_arena_alloc(state, text.len + 1, 1, &storage);
    if (err.msg.data) return err;
    unsigned char *data = storage;
    if (text.len) memcpy(data, text.data, text.len);
    data[text.len] = 0;
    *result = (_String){data, text.len};
    return (_Error){0};
}
