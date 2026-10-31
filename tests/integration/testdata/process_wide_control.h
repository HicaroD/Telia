#include "allocation_control.h"
#include <windows.h>
#include <wchar.h>

static int _test_wide_case;
static void _test_wide_mode(int32_t mode) { _test_wide_case = mode; }
static void _test_native_wide_args(int *argc, wchar_t ***argv) {
    static wchar_t invalid[] = {0xd800, 0};
    static wchar_t valid[] = L"valid";
    static wchar_t *values[] = {L"program", valid, invalid};
    *argc = 3;
    *argv = values;
}
static int _test_wide_conversion(UINT page, DWORD flags, const wchar_t *text,
        int len, char *out, int size, const char *fallback, BOOL *used) {
    // Mode 1 uses an actual malformed surrogate; mode 2/3 inject conversion
    // failure in the sizing/output call after a prior argument was copied.
    if (text[0] == 0xd800 && ((_test_wide_case == 2 && !out) ||
                             (_test_wide_case == 3 && out))) return 0;
    if (text[0] == 0xd800 && _test_wide_case != 1)
        return WideCharToMultiByte(page, flags, L"x", 1, out, size, fallback, used);
    return WideCharToMultiByte(page, flags, text, len, out, size, fallback, used);
}
#define _TELIA_NATIVE_WIDE_ARGS _test_native_wide_args
#define _TELIA_WIDE_TO_UTF8 _test_wide_conversion
