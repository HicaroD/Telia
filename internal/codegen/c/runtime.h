// Private implementation of the standard library's raw-memory operations.
// Public contracts live in base/std; this file uses only the primitive ABI.
#define _TELIA_ERROR(kind_name, message) ((_Error){.msg = {(const unsigned char *)(message), sizeof(message) - 1}, .present = true, .kind = &_telia_runtime_kind_##kind_name})

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

static inline _Error _telia_string_byte(_String text, uint64_t index, uint32_t *result) {
    *result = 0;
    if (index >= text.len) return _TELIA_ERROR(OutOfBounds, "string byte index out of bounds");
    *result = text.data[(size_t)index];
    return (_Error){0};
}

static inline _Error _telia_string_slice(_String text, uint64_t start,
                                        uint64_t end, _String *result) {
    *result = (_String){0};
    if (start > end || end > text.len) {
        return _TELIA_ERROR(InvalidRange, "invalid string range");
    }
    if (start != end) {
        *result = (_String){text.data + (size_t)start, (size_t)(end - start)};
    }
    return (_Error){0};
}

static inline _Error _telia_string_concat(_String left, _String right, _String *result) {
    *result = (_String){0};
    if (right.len > SIZE_MAX - left.len) return _TELIA_ERROR(SizeOverflow, "allocation size overflow");
    size_t size = left.len + right.len;
    if (!size) return (_Error){0};
    unsigned char *data = _TELIA_MALLOC(size);
    if (!data) return _TELIA_ERROR(AllocationFailed, "allocation failed");
    if (left.len) memcpy(data, left.data, left.len);
    if (right.len) memcpy(data + left.len, right.data, right.len);
    *result = (_String){data, size};
    return (_Error){0};
}

static inline _Error _telia_string_clone(_String text, _String *result) {
    return _telia_string_concat(text, (_String){0}, result);
}

// Bounded output construction only; parsing and signed magnitude handling live
// in std::strconv. Twenty decimal digits plus an optional minus sign suffice.
static inline _Error _telia_format_decimal(uint64_t magnitude, bool negative, _String *result) {
    unsigned char buffer[21];
    size_t start = sizeof(buffer);
    do {
        buffer[--start] = (unsigned char)('0' + magnitude % 10);
        magnitude /= 10;
    } while (magnitude);
    if (negative) buffer[--start] = '-';
    return _telia_string_clone((_String){buffer + start, sizeof(buffer) - start}, result);
}

static inline void _telia_string_free(_String *text) {
    _TELIA_FREE((void *)text->data);
    *text = (_String){0};
}

static inline _Error _telia_string_to_c(_String text, _String *result) {
    *result = (_String){0};
    if (text.len == SIZE_MAX) return _TELIA_ERROR(SizeOverflow, "allocation size overflow");
    if (text.len && memchr(text.data, 0, text.len)) return _TELIA_ERROR(EmbeddedZero, "embedded zero byte");
    unsigned char *data = _TELIA_MALLOC(text.len + 1);
    if (!data) return _TELIA_ERROR(AllocationFailed, "allocation failed");
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
    if (size > SIZE_MAX) return _TELIA_ERROR(SizeOverflow, "allocation size overflow");
    if (!size) return (_Error){0};
    void *data = _TELIA_MALLOC((size_t)size);
    if (!data) return _TELIA_ERROR(AllocationFailed, "allocation failed");
    *result = data;
    return (_Error){0};
}

static inline void _telia_memory_free(void *ptr) {
    _TELIA_FREE(ptr);
}

#ifdef _WIN32
#include <windows.h>
#include <wchar.h>
#endif

// Private OS boundary overrides are used only by native failure fixtures.
#ifndef _TELIA_FS_OPEN
#ifdef _WIN32
#define _TELIA_FS_OPEN _wfopen
#else
#define _TELIA_FS_OPEN fopen
#endif
#endif
#ifndef _TELIA_FS_READ
#define _TELIA_FS_READ fread
#endif
#ifndef _TELIA_FS_WRITE
#define _TELIA_FS_WRITE fwrite
#endif
#ifndef _TELIA_FS_ERROR
#define _TELIA_FS_ERROR ferror
#endif
#ifndef _TELIA_FS_EOF
#define _TELIA_FS_EOF feof
#endif
#ifndef _TELIA_FS_CLOSE
#define _TELIA_FS_CLOSE fclose
#endif

static inline _Error _telia_fs_open(_String path, bool writing, FILE **result) {
    *result = NULL;
#ifdef _WIN32
    if (path.len > INT_MAX) return _TELIA_ERROR(SizeOverflow, "path too large");
    int count = MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS,
        (const char *)path.data, (int)path.len, NULL, 0);
    if (!count) return _TELIA_ERROR(InvalidUTF8, "invalid UTF-8 path");
    if ((size_t)count > SIZE_MAX / sizeof(wchar_t) - 1)
        return _TELIA_ERROR(SizeOverflow, "path too large");
    wchar_t *native = _TELIA_MALLOC(((size_t)count + 1) * sizeof(wchar_t));
    if (!native) return _TELIA_ERROR(AllocationFailed, "allocation failed");
    if (!MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS,
        (const char *)path.data, (int)path.len, native, count)) {
        _TELIA_FREE(native);
        return _TELIA_ERROR(InvalidUTF8, "invalid UTF-8 path");
    }
    native[count] = 0;
    *result = _TELIA_FS_OPEN(native, writing ? L"wb" : L"rb");
    _TELIA_FREE(native);
#else
    _String native = {0};
    _Error err = _telia_string_to_c(path, &native);
    if (err.present) return err;
    *result = _TELIA_FS_OPEN((const char *)native.data, writing ? "wb" : "rb");
    _telia_string_free(&native);
#endif
    if (!*result) return _TELIA_ERROR(FileOpenFailed, "file open failed");
    return (_Error){0};
}

static inline _Error _telia_fs_read(_String path, _String *result) {
    *result = (_String){0};
    FILE *file = NULL;
    _Error err = _telia_fs_open(path, false, &file);
    if (err.present) return err;
    unsigned char *data = NULL;
    size_t len = 0, capacity = 0;
    unsigned char chunk[8192];
    for (;;) {
        size_t count = _TELIA_FS_READ(chunk, 1, sizeof(chunk), file);
        if (count) {
            if (count > SIZE_MAX - len) {
                err = _TELIA_ERROR(SizeOverflow, "file too large");
                break;
            }
            size_t needed = len + count;
            if (needed > capacity) {
                size_t next = capacity ? capacity : sizeof(chunk);
                while (next < needed) {
                    if (next > SIZE_MAX / 2) { next = needed; break; }
                    next *= 2;
                }
                unsigned char *grown = _TELIA_MALLOC(next);
                if (!grown) { err = _TELIA_ERROR(AllocationFailed, "allocation failed"); break; }
                if (len) memcpy(grown, data, len);
                _TELIA_FREE(data);
                data = grown;
                capacity = next;
            }
            memcpy(data + len, chunk, count);
            len += count;
        }
        if (_TELIA_FS_ERROR(file)) { err = _TELIA_ERROR(FileReadFailed, "file read failed"); break; }
        if (_TELIA_FS_EOF(file)) break;
        if (!count) { err = _TELIA_ERROR(FileReadFailed, "file read made no progress"); break; }
    }
    if (_TELIA_FS_CLOSE(file) != 0 && !err.present)
        err = _TELIA_ERROR(FileCloseFailed, "file close failed");
    if (err.present) _TELIA_FREE(data);
    else *result = (_String){data, len};
    return err;
}

static inline _Error _telia_fs_write(_String path, _String data) {
    FILE *file = NULL;
    _Error err = _telia_fs_open(path, true, &file);
    if (err.present) return err;
    size_t offset = 0;
    while (offset < data.len) {
        size_t count = _TELIA_FS_WRITE(data.data + offset, 1, data.len - offset, file);
        offset += count;
        if (_TELIA_FS_ERROR(file) || !count) {
            err = _TELIA_ERROR(FileWriteFailed, "file write failed");
            break;
        }
    }
    if (_TELIA_FS_CLOSE(file) != 0 && !err.present)
        err = _TELIA_ERROR(FileCloseFailed, "file close failed");
    return err;
}

// Standard streams remain owned by the process. Overrides are native test seams.
#ifndef _TELIA_IO_MAX_SIZE
#define _TELIA_IO_MAX_SIZE SIZE_MAX
#endif
#ifndef _TELIA_IO_READ
#define _TELIA_IO_READ fread
#endif
#ifndef _TELIA_IO_EOF
#define _TELIA_IO_EOF feof
#endif
#ifndef _TELIA_IO_WRITE
#define _TELIA_IO_WRITE fwrite
#endif
#ifndef _TELIA_IO_FLUSH
#define _TELIA_IO_FLUSH fflush
#endif
#ifndef _TELIA_IO_ERROR
#define _TELIA_IO_ERROR ferror
#endif
#ifdef _WIN32
#ifndef _TELIA_IO_SETMODE
#define _TELIA_IO_SETMODE _setmode
#endif
#else
#include <signal.h>
#include <pthread.h>
#ifdef __linux__
#include <stdio_ext.h>
#endif
#endif

static inline _Error _telia_io_read(_String *result) {
    *result = (_String){0};
    _Error err = {0};
#ifdef _WIN32
    int mode = _TELIA_IO_SETMODE(_fileno(stdin), _O_BINARY);
    if (mode == -1) return _TELIA_ERROR(StreamModeFailed, "stream mode failed");
#endif
    unsigned char *data = NULL;
    size_t len = 0, capacity = 0;
    unsigned char chunk[8192];
    for (;;) {
        size_t count = _TELIA_IO_READ(chunk, 1, sizeof(chunk), stdin);
        if (count) {
            if (count > _TELIA_IO_MAX_SIZE - len) {
                err = _TELIA_ERROR(SizeOverflow, "input too large");
                break;
            }
            size_t needed = len + count;
            if (needed > capacity) {
                size_t next = capacity ? capacity : sizeof(chunk);
                while (next < needed) {
                    if (next > SIZE_MAX / 2) { next = needed; break; }
                    next *= 2;
                }
                unsigned char *grown = _TELIA_MALLOC(next);
                if (!grown) { err = _TELIA_ERROR(AllocationFailed, "allocation failed"); break; }
                if (len) memcpy(grown, data, len);
                _TELIA_FREE(data);
                data = grown;
                capacity = next;
            }
            memcpy(data + len, chunk, count);
            len += count;
        }
        if (_TELIA_IO_ERROR(stdin)) { err = _TELIA_ERROR(StreamReadFailed, "stream read failed"); break; }
        if (_TELIA_IO_EOF(stdin)) break;
        if (!count) { err = _TELIA_ERROR(StreamReadFailed, "stream read made no progress"); break; }
    }
#ifdef _WIN32
    if (_TELIA_IO_SETMODE(_fileno(stdin), mode) == -1 && !err.present)
        err = _TELIA_ERROR(StreamModeFailed, "stream mode restore failed");
#endif
    if (err.present) _TELIA_FREE(data);
    else *result = (_String){data, len};
    return err;
}

static inline _Error _telia_io_write(_String data, bool to_stderr) {
    FILE *stream = to_stderr ? stderr : stdout;
    _Error err = {0};
#ifdef _WIN32
    // Flush existing text output before changing its descriptor's translation.
    if (_TELIA_IO_FLUSH(stream) != 0)
        return _TELIA_ERROR(StreamFlushFailed, "stream flush failed");
    int mode = _TELIA_IO_SETMODE(_fileno(stream), _O_BINARY);
    if (mode == -1) return _TELIA_ERROR(StreamModeFailed, "stream mode failed");
#else
    // Block only this thread's SIGPIPE, without changing the process handler.
    // Consume a newly generated signal before restoring the caller's mask.
    sigset_t blocked, previous, pending;
    sigemptyset(&blocked);
    sigaddset(&blocked, SIGPIPE);
    if (pthread_sigmask(SIG_BLOCK, &blocked, &previous) != 0)
        return _TELIA_ERROR(StreamWriteFailed, "cannot block SIGPIPE");
    if (sigpending(&pending) != 0) {
        pthread_sigmask(SIG_SETMASK, &previous, NULL);
        return _TELIA_ERROR(StreamWriteFailed, "cannot inspect pending signals");
    }
    bool already_pending = sigismember(&pending, SIGPIPE) == 1;
#endif
    size_t offset = 0;
    while (offset < data.len) {
        size_t count = _TELIA_IO_WRITE(data.data + offset, 1, data.len - offset, stream);
        offset += count;
        if (_TELIA_IO_ERROR(stream) || !count) {
            err = _TELIA_ERROR(StreamWriteFailed, "stream write failed");
            break;
        }
    }
    if (_TELIA_IO_FLUSH(stream) != 0 && !err.present)
        err = _TELIA_ERROR(StreamFlushFailed, "stream flush failed");
#ifdef _WIN32
    if (_TELIA_IO_SETMODE(_fileno(stream), mode) == -1 && !err.present)
        err = _TELIA_ERROR(StreamModeFailed, "stream mode restore failed");
#else
    // Some libcs retain a failed flush's buffer and retry it at process exit.
    // Discard unsent bytes so a checked broken pipe cannot kill the caller later.
    if (err.present) {
#ifdef __linux__
        __fpurge(stream);
#else
        fpurge(stream);
#endif
    }
    if (!already_pending && sigpending(&pending) == 0 && sigismember(&pending, SIGPIPE) == 1) {
        int signal_number;
        sigwait(&blocked, &signal_number);
    }
    if (pthread_sigmask(SIG_SETMASK, &previous, NULL) != 0 && !err.present)
        err = _TELIA_ERROR(StreamWriteFailed, "cannot restore signal mask");
#endif
    return err;
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
        return _TELIA_ERROR(InvalidAlignment, "invalid allocation alignment");
    }
    if (size > SIZE_MAX || size > SIZE_MAX - (alignment - 1)) {
        return _TELIA_ERROR(SizeOverflow, "allocation size overflow");
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
        return _TELIA_ERROR(SizeOverflow, "allocation size overflow");
    }
    _TeliaArenaBlock *block = _TELIA_MALLOC(sizeof(*block) + capacity);
    if (!block) return _TELIA_ERROR(AllocationFailed, "allocation failed");
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
    if (right.len > SIZE_MAX - left.len) return _TELIA_ERROR(SizeOverflow, "allocation size overflow");
    size_t size = left.len + right.len;
    void *storage;
    _Error err = _telia_arena_alloc(state, size, 1, &storage);
    if (err.present) return err;
    unsigned char *data = storage;
    if (left.len) memcpy(data, left.data, left.len);
    if (right.len) memcpy(data + left.len, right.data, right.len);
    *result = (_String){data, size};
    return (_Error){0};
}

static inline _Error _telia_arena_string_to_c(void **state, _String text, _String *result) {
    *result = (_String){0};
    if (text.len == SIZE_MAX) return _TELIA_ERROR(SizeOverflow, "allocation size overflow");
    if (text.len && memchr(text.data, 0, text.len)) return _TELIA_ERROR(EmbeddedZero, "embedded zero byte");
    void *storage;
    _Error err = _telia_arena_alloc(state, text.len + 1, 1, &storage);
    if (err.present) return err;
    unsigned char *data = storage;
    if (text.len) memcpy(data, text.data, text.len);
    data[text.len] = 0;
    *result = (_String){data, text.len};
    return (_Error){0};
}
