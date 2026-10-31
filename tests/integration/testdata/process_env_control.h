#include "allocation_control.h"
#ifdef _WIN32
#include <windows.h>
static void _test_change_environment(void) {
    if (!SetEnvironmentVariableW(L"TELIA_ENV_TEST", L"changed")) abort();
}
#else
#include <string.h>
static void _test_change_environment(void) {
    char *value = getenv("TELIA_ENV_TEST");
    if (!value || strcmp(value, "original")) abort();
    // Mutate the actual native buffer, so borrowing it cannot pass this test.
    memcpy(value, "changed!", 8);
}
#endif
