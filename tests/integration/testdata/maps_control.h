#include "allocation_control.h"
#include <string.h>
#include <math.h>

static int64_t map_seed_calls = 0;
static bool map_seed_fail = false;
static int map_entropy(void *bytes, size_t size) {
    ++map_seed_calls;
    memset(bytes, 0x71, size);
    return map_seed_fail ? -1 : 0;
}
static void _map_seed_failure(bool fail) { map_seed_fail = fail; }
static int64_t _map_seed_calls(void) { return map_seed_calls; }
static double _map_nan(void) { return NAN; }
static double _map_negative_zero(void) { return -0.0; }
static const char *_map_malformed(void) { return "\xff\xfe"; }
#define _TELIA_ENTROPY(bytes, size) map_entropy(bytes, size)
#ifdef TELIA_MAP_COLLISIONS
// Identical fingerprints, final group start, wraparound and long chains.
#define _TELIA_KEY_HASH(value) ((void)(value), UINT64_MAX)
#endif
