#include "allocation_control.h"
#include <windows.h>
#include <wchar.h>

static int32_t _test_env_case;
static void _test_env_mode(int32_t mode) { _test_env_case = mode; }
static DWORD WINAPI _test_get_environment(LPCWSTR name, LPWSTR out, DWORD size) {
    (void)name;
    if (_test_env_case == 4) {
        SetLastError(ERROR_ACCESS_DENIED);
        return 0;
    }
    if (!out) return 2;
    if (_test_env_case == 5) return size + 1;
    out[0] = _test_env_case == 1 ? 0xd800 : L'x';
    out[1] = 0;
    return 1;
}
static int _test_env_conversion(UINT page, DWORD flags, const wchar_t *text,
        int len, char *out, int size, const char *fallback, BOOL *used) {
    if ((_test_env_case == 2 && !out) || (_test_env_case == 3 && out)) return 0;
    return WideCharToMultiByte(page, flags, text, len, out, size, fallback, used);
}
static int _test_env_name_conversion(UINT page, DWORD flags, const char *text,
        int len, wchar_t *out, int size) {
    if ((_test_env_case == 6 && !out) || (_test_env_case == 7 && out)) return 0;
    return MultiByteToWideChar(page, flags, text, len, out, size);
}
#define _TELIA_GET_ENVIRONMENT _test_get_environment
#define _TELIA_WIDE_TO_UTF8 _test_env_conversion
#define _TELIA_ENV_NAME_TO_WIDE _test_env_name_conversion
