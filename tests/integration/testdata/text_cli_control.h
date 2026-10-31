#include "io_control.h"

// Test-only injection at the native boundary; the example has no C bindings.
__attribute__((constructor)) static void _test_text_cli_controls(void) {
    const char *limit = getenv("TELIA_TEXT_FAIL_AFTER");
    if (limit) _test_fail_after(strtoll(limit, NULL, 10));
    const char *mode = getenv("TELIA_TEXT_IO_MODE");
    if (mode) _test_io_mode((int32_t)strtol(mode, NULL, 10));
}
