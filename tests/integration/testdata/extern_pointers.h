#include <stdint.h>

static int32_t *telia_test_pointer_round_trip(int32_t **slot) {
    **slot += 2;
    return *slot;
}
