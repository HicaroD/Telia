#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

static int64_t test_remaining = -1;
static int64_t test_live = 0;

static void *_test_malloc(size_t size) {
    if (test_remaining == 0) return NULL;
    if (test_remaining > 0) --test_remaining;
    void *ptr = malloc(size);
    if (ptr) ++test_live;
    return ptr;
}

static void _test_free(void *ptr) {
    if (ptr) --test_live;
    free(ptr);
}

static void _test_fail_after(int64_t count) { test_remaining = count; }
static int64_t _test_live_allocations(void) { return test_live; }
static bool _test_is_aligned(void *ptr, uint64_t alignment) {
    return (uintptr_t)ptr % alignment == 0;
}

__attribute__((destructor)) static void _test_require_no_leaks(void) {
    if (test_live != 0) abort();
}

#define _TELIA_MALLOC _test_malloc
#define _TELIA_FREE _test_free
