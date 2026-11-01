#pragma once
#include <stddef.h>
static int hash_test_entropy(unsigned char *bytes, size_t n) {
    for (size_t i = 0; i < n; ++i) bytes[i] = (unsigned char)i;
    return TELIA_ENTROPY_FAIL ? -1 : 0; // Failure can leave partially written bytes.
}
static int hash_entropy_fails(void) { return TELIA_ENTROPY_FAIL; }
#define _TELIA_ENTROPY hash_test_entropy
