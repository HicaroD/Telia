#ifndef TELIA_TEST_TERMINAL_CONSOLE_H
#define TELIA_TEST_TERMINAL_CONSOLE_H
#include <windows.h>
#include <stdint.h>

static DWORD terminal_test_input_mode, terminal_test_output_mode;
static UINT terminal_test_input_cp, terminal_test_output_cp;

static inline int32_t terminal_test_console(void) {
    // This test process owns its console; never modify the test runner's console.
    FreeConsole();
    if (!AllocConsole()) return 1;
    HANDLE input = CreateFileW(L"CONIN$", GENERIC_READ | GENERIC_WRITE,
                               FILE_SHARE_READ | FILE_SHARE_WRITE, NULL, OPEN_EXISTING, 0, NULL);
    HANDLE output = CreateFileW(L"CONOUT$", GENERIC_READ | GENERIC_WRITE,
                                FILE_SHARE_READ | FILE_SHARE_WRITE, NULL, OPEN_EXISTING, 0, NULL);
    if (input == INVALID_HANDLE_VALUE || output == INVALID_HANDLE_VALUE) return 2;
    BOOL ok = GetConsoleMode(input, &terminal_test_input_mode) &&
              GetConsoleMode(output, &terminal_test_output_mode);
    terminal_test_input_cp = GetConsoleCP();
    terminal_test_output_cp = GetConsoleOutputCP();
    CloseHandle(input);
    CloseHandle(output);
    return ok ? 0 : 3;
}

static inline int32_t terminal_test_key(intptr_t handle) {
    INPUT_RECORD record = {0};
    record.EventType = KEY_EVENT;
    record.Event.KeyEvent.bKeyDown = TRUE;
    record.Event.KeyEvent.wRepeatCount = 1;
    record.Event.KeyEvent.wVirtualKeyCode = 'E';
    record.Event.KeyEvent.uChar.UnicodeChar = 0x00e9;
    DWORD written;
    return WriteConsoleInputW((HANDLE)handle, &record, 1, &written) && written == 1 ? 0 : 1;
}

static inline int32_t terminal_test_scalar(intptr_t handle, uint32_t scalar) {
    INPUT_RECORD records[2] = {0};
    DWORD count = scalar > 0xffff ? 2 : 1;
    for (DWORD i = 0; i < count; ++i) {
        records[i].EventType = KEY_EVENT;
        records[i].Event.KeyEvent.bKeyDown = TRUE;
        records[i].Event.KeyEvent.wRepeatCount = 1;
    }
    if (count == 1) {
        records[0].Event.KeyEvent.uChar.UnicodeChar = (WCHAR)scalar;
    } else {
        scalar -= 0x10000;
        records[0].Event.KeyEvent.uChar.UnicodeChar = (WCHAR)(0xd800 + (scalar >> 10));
        records[1].Event.KeyEvent.uChar.UnicodeChar = (WCHAR)(0xdc00 + (scalar & 0x3ff));
    }
    DWORD written;
    return WriteConsoleInputW((HANDLE)handle, records, count, &written) && written == count ? 0 : 1;
}

static inline int32_t terminal_test_resize(intptr_t output) {
    CONSOLE_SCREEN_BUFFER_INFO info;
    if (!GetConsoleScreenBufferInfo((HANDLE)output, &info)) return 1;
    SMALL_RECT window = info.srWindow;
    --window.Right;
    --window.Bottom;
    return SetConsoleWindowInfo((HANDLE)output, TRUE, &window) ? 0 : 2;
}

static inline int32_t terminal_test_restored(intptr_t input, intptr_t output) {
    DWORD input_mode, output_mode;
    return GetConsoleMode((HANDLE)input, &input_mode) &&
           GetConsoleMode((HANDLE)output, &output_mode) &&
           input_mode == terminal_test_input_mode && output_mode == terminal_test_output_mode &&
           GetConsoleCP() == terminal_test_input_cp && GetConsoleOutputCP() == terminal_test_output_cp ? 0 : 1;
}
#endif
