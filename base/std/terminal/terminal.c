#ifndef _WIN32
#define _POSIX_C_SOURCE 200809L
#endif
#include "telia_terminal.h"
#include <limits.h>
#include <stdlib.h>
#include <time.h>
#include <string.h>
#ifdef _WIN32
#include <windows.h>
#else
#include <errno.h>
#include <unistd.h>
#include <termios.h>
#include <sys/ioctl.h>
#include <poll.h>
#include <fcntl.h>
#include <signal.h>
#include <pthread.h>
#endif

static bool terminal_handle_in_range(intptr_t handle) {
#ifdef _WIN32
    return handle != 0 && (HANDLE)handle != INVALID_HANDLE_VALUE;
#else
    return handle >= 0 && handle <= INT_MAX;
#endif
}

intptr_t telia_terminal_standard(int32_t stream) {
#ifdef _WIN32
    return (intptr_t)GetStdHandle(stream == 0 ? STD_INPUT_HANDLE :
                                 stream == 1 ? STD_OUTPUT_HANDLE : STD_ERROR_HANDLE);
#else
    return stream;
#endif
}

int32_t telia_terminal_is_terminal(intptr_t handle, bool *result) {
    *result = false;
    if (!terminal_handle_in_range(handle)) return 6;
#ifdef _WIN32
    DWORD mode;
    *result = GetConsoleMode((HANDLE)handle, &mode) != 0;
    if (!*result) {
        SetLastError(NO_ERROR);
        if (GetFileType((HANDLE)handle) == FILE_TYPE_UNKNOWN && GetLastError() != NO_ERROR) return 6;
    }
#else
    *result = isatty((int)handle) != 0;
    if (!*result && errno != ENOTTY) return 6;
#endif
    return 0;
}

int32_t telia_terminal_size(intptr_t handle, uint64_t *columns, uint64_t *rows) {
    *columns = *rows = 0;
    bool attached;
    int32_t status = telia_terminal_is_terminal(handle, &attached);
    if (status) return status;
    if (!attached) return 4;
#ifdef _WIN32
    CONSOLE_SCREEN_BUFFER_INFO info;
    if (!GetConsoleScreenBufferInfo((HANDLE)handle, &info)) {
        HANDLE output = CreateFileW(L"CONOUT$", GENERIC_READ | GENERIC_WRITE,
                                    FILE_SHARE_READ | FILE_SHARE_WRITE, NULL, OPEN_EXISTING, 0, NULL);
        if (output == INVALID_HANDLE_VALUE) return 6;
        BOOL ok = GetConsoleScreenBufferInfo(output, &info);
        CloseHandle(output);
        if (!ok) return 6;
    }
    *columns = (uint64_t)(info.srWindow.Right - info.srWindow.Left + 1);
    *rows = (uint64_t)(info.srWindow.Bottom - info.srWindow.Top + 1);
#else
    struct winsize dimensions;
    if (ioctl((int)handle, TIOCGWINSZ, &dimensions) != 0) return 6;
    *columns = dimensions.ws_col;
    *rows = dimensions.ws_row;
#endif
    return 0;
}

struct terminal_saved_mode {
    intptr_t handle;
#ifdef _WIN32
    DWORD mode;
    UINT codepage;
    bool output;
#else
    struct termios mode;
#endif
};

int32_t telia_terminal_get_mode(intptr_t handle, bool *echo, bool *canonical, bool *signals) {
    bool attached;
    int32_t status = telia_terminal_is_terminal(handle, &attached);
    if (status) return status;
    if (!attached) return 4;
#ifdef _WIN32
    DWORD mode, events;
    if (!GetNumberOfConsoleInputEvents((HANDLE)handle, &events)) return 5;
    if (!GetConsoleMode((HANDLE)handle, &mode)) return 6;
    *echo = (mode & ENABLE_ECHO_INPUT) != 0;
    *canonical = (mode & ENABLE_LINE_INPUT) != 0;
    *signals = (mode & ENABLE_PROCESSED_INPUT) != 0;
#else
    struct termios mode;
    if (tcgetattr((int)handle, &mode) != 0) return 6;
    *echo = (mode.c_lflag & ECHO) != 0;
    *canonical = (mode.c_lflag & ICANON) != 0;
    *signals = (mode.c_lflag & ISIG) != 0;
#endif
    return 0;
}

int32_t telia_terminal_apply_mode(intptr_t handle, bool raw, bool echo, bool canonical, bool signals, void **saved) {
    *saved = NULL;
    bool a, b, c;
    int32_t status = telia_terminal_get_mode(handle, &a, &b, &c);
    if (status) return status;
    struct terminal_saved_mode *snapshot = malloc(sizeof(*snapshot));
    if (!snapshot) return 6;
    snapshot->handle = handle;
#ifdef _WIN32
    snapshot->output = false;
    snapshot->codepage = GetConsoleCP();
    if (!GetConsoleMode((HANDLE)handle, &snapshot->mode)) { free(snapshot); return 6; }
    DWORD mode = snapshot->mode;
    if (raw) mode &= ~(ENABLE_LINE_INPUT | ENABLE_ECHO_INPUT | ENABLE_PROCESSED_INPUT | ENABLE_QUICK_EDIT_MODE);
    mode |= ENABLE_VIRTUAL_TERMINAL_INPUT | ENABLE_EXTENDED_FLAGS | ENABLE_WINDOW_INPUT;
    mode = echo ? mode | ENABLE_ECHO_INPUT : mode & ~ENABLE_ECHO_INPUT;
    mode = canonical ? mode | ENABLE_LINE_INPUT : mode & ~ENABLE_LINE_INPUT;
    mode = signals ? mode | ENABLE_PROCESSED_INPUT : mode & ~ENABLE_PROCESSED_INPUT;
    if (!SetConsoleMode((HANDLE)handle, mode)) { free(snapshot); return 5; }
    if (!SetConsoleCP(CP_UTF8)) {
        SetConsoleMode((HANDLE)handle, snapshot->mode);
        free(snapshot);
        return 6;
    }
#else
    if (tcgetattr((int)handle, &snapshot->mode) != 0) { free(snapshot); return 6; }
    struct termios mode = snapshot->mode;
    if (raw) {
        mode.c_iflag &= ~(IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | IGNCR | ICRNL | IXON);
        mode.c_oflag &= ~OPOST;
        mode.c_lflag &= ~(ECHO | ECHONL | ICANON | ISIG | IEXTEN);
        mode.c_cflag = (mode.c_cflag & ~(CSIZE | PARENB)) | CS8;
    }
    mode.c_lflag = echo ? mode.c_lflag | ECHO : mode.c_lflag & ~ECHO;
    mode.c_lflag = canonical ? mode.c_lflag | ICANON : mode.c_lflag & ~ICANON;
    mode.c_lflag = signals ? mode.c_lflag | ISIG : mode.c_lflag & ~ISIG;
    if (!canonical) { mode.c_cc[VMIN] = 1; mode.c_cc[VTIME] = 0; }
    if (tcsetattr((int)handle, TCSANOW, &mode) != 0) { free(snapshot); return 6; }
#endif
    *saved = snapshot;
    return 0;
}

int32_t telia_terminal_restore(void **saved) {
    struct terminal_saved_mode *snapshot = *saved;
    if (!snapshot) return 0;
#ifdef _WIN32
    if (!SetConsoleMode((HANDLE)snapshot->handle, snapshot->mode)) return 6;
    if (snapshot->output ? !SetConsoleOutputCP(snapshot->codepage) : !SetConsoleCP(snapshot->codepage)) return 6;
#else
    if (tcsetattr((int)snapshot->handle, TCSANOW, &snapshot->mode) != 0) return 6;
#endif
    free(snapshot);
    *saved = NULL;
    return 0;
}

int32_t telia_terminal_prepare_output(intptr_t handle, void **saved) {
    *saved = NULL;
    bool attached;
    int32_t status = telia_terminal_is_terminal(handle, &attached);
    if (status) return status;
    if (!attached) return 4;
#ifdef _WIN32
    struct terminal_saved_mode *snapshot = malloc(sizeof(*snapshot));
    if (!snapshot) return 6;
    snapshot->handle = handle;
    snapshot->output = true;
    snapshot->codepage = GetConsoleOutputCP();
    if (!GetConsoleMode((HANDLE)handle, &snapshot->mode)) { free(snapshot); return 6; }
    if (!SetConsoleMode((HANDLE)handle, snapshot->mode | ENABLE_PROCESSED_OUTPUT | ENABLE_VIRTUAL_TERMINAL_PROCESSING)) {
        free(snapshot);
        return 5;
    }
    if (!SetConsoleOutputCP(CP_UTF8)) {
        SetConsoleMode((HANDLE)handle, snapshot->mode);
        free(snapshot);
        return 6;
    }
    *saved = snapshot;
#endif
    return 0;
}

int32_t telia_terminal_profile(intptr_t handle, bool *dumb, bool *ansi, bool *indexed, bool *rgb) {
    *dumb = *ansi = *indexed = *rgb = false;
    bool attached;
    int32_t status = telia_terminal_is_terminal(handle, &attached);
    if (status) return status;
    if (!attached) { *dumb = true; return 0; }
    const char *term = getenv("TERM");
    const char *color = getenv("COLORTERM");
    *dumb = term && strcmp(term, "dumb") == 0;
    if (*dumb) return 0;
    if (term) {
        *ansi = strncmp(term, "xterm", 5) == 0 || strncmp(term, "screen", 6) == 0 ||
                strncmp(term, "tmux", 4) == 0 || strcmp(term, "vt100") == 0 ||
                strcmp(term, "ansi") == 0 || strncmp(term, "linux", 5) == 0;
        *indexed = strstr(term, "256color") != NULL;
        *rgb = strstr(term, "direct") != NULL;
    }
    if (color && (strcmp(color, "truecolor") == 0 || strcmp(color, "24bit") == 0)) *rgb = true;
    if (*rgb || *indexed) *ansi = true;
    return 0;
}

int32_t telia_terminal_close(intptr_t handle) {
    if (!terminal_handle_in_range(handle)) return 6;
#ifdef _WIN32
    return CloseHandle((HANDLE)handle) ? 0 : 6;
#else
    return close((int)handle) == 0 ? 0 : 6;
#endif
}

int32_t telia_terminal_open(bool input, intptr_t *handle) {
#ifdef _WIN32
    HANDLE opened = CreateFileW(input ? L"CONIN$" : L"CONOUT$", GENERIC_READ | GENERIC_WRITE,
                                FILE_SHARE_READ | FILE_SHARE_WRITE, NULL, OPEN_EXISTING, 0, NULL);
    if (opened == INVALID_HANDLE_VALUE) return 6;
    *handle = (intptr_t)opened;
#else
    (void)input;
    int opened = open("/dev/tty", O_RDWR | O_CLOEXEC);
    if (opened < 0) return 6;
    *handle = opened;
#endif
    return 0;
}

int32_t telia_terminal_write(intptr_t handle, void *bytes, uint64_t length, uint64_t *count) {
    *count = 0;
    if (!terminal_handle_in_range(handle)) return 6;
    if (length == 0) return 0;
#ifdef _WIN32
    DWORD written;
    DWORD chunk = length > MAXDWORD ? MAXDWORD : (DWORD)length;
    if (!WriteFile((HANDLE)handle, bytes, chunk, &written, NULL)) return 6;
    *count = written;
#else
    size_t chunk = length > (uint64_t)SSIZE_MAX ? (size_t)SSIZE_MAX : (size_t)length;
    // Suppress only the SIGPIPE generated by this write, without installing a
    // process-wide handler or consuming a signal that was already pending.
    sigset_t blocked, previous, pending;
    sigemptyset(&blocked);
    sigaddset(&blocked, SIGPIPE);
    if (pthread_sigmask(SIG_BLOCK, &blocked, &previous) != 0) return 6;
    if (sigpending(&pending) != 0) {
        pthread_sigmask(SIG_SETMASK, &previous, NULL);
        return 6;
    }
    bool already_pending = sigismember(&pending, SIGPIPE) == 1;
    ssize_t written;
    do { written = write((int)handle, bytes, chunk); } while (written < 0 && errno == EINTR);
    int failure = errno;
    if (written < 0 && failure == EPIPE && !already_pending) {
        if (sigpending(&pending) == 0 && sigismember(&pending, SIGPIPE) == 1) {
            int signal_number;
            sigwait(&blocked, &signal_number);
        }
    }
    pthread_sigmask(SIG_SETMASK, &previous, NULL);
    errno = failure;
    if (written < 0) return 6;
    *count = (uint64_t)written;
#endif
    return 0;
}

int32_t telia_terminal_read(intptr_t handle, void *bytes, uint64_t length, uint64_t *count, bool *eof) {
    *count = 0;
    *eof = false;
    if (!terminal_handle_in_range(handle)) return 6;
    if (length == 0) return 0;
#ifdef _WIN32
    DWORD received;
    DWORD chunk = length > MAXDWORD ? MAXDWORD : (DWORD)length;
    if (!ReadFile((HANDLE)handle, bytes, chunk, &received, NULL)) {
        if (GetLastError() == ERROR_BROKEN_PIPE || GetLastError() == ERROR_HANDLE_EOF) {
            *eof = true;
            return 0;
        }
        return 6;
    }
    *count = received;
#else
    size_t chunk = length > (uint64_t)SSIZE_MAX ? (size_t)SSIZE_MAX : (size_t)length;
    ssize_t received;
    do { received = read((int)handle, bytes, chunk); } while (received < 0 && errno == EINTR);
    if (received < 0) return 6;
    *count = (uint64_t)received;
#endif
    *eof = *count == 0;
    return 0;
}

int32_t telia_terminal_write_text(intptr_t handle, char *bytes, uint64_t length, uint64_t *count) {
    return telia_terminal_write(handle, bytes, length, count);
}

static uint64_t terminal_milliseconds(void) {
#ifdef _WIN32
    return GetTickCount64();
#else
    struct timespec now;
    clock_gettime(CLOCK_MONOTONIC, &now);
    return (uint64_t)now.tv_sec * 1000 + (uint64_t)now.tv_nsec / 1000000;
#endif
}

int32_t telia_terminal_wait(intptr_t handle, bool input, bool resize, int64_t timeout_ms,
                           bool *initialized, uint64_t *columns, uint64_t *rows,
                           bool *ready, bool *resized, bool *timed_out) {
    *ready = *resized = *timed_out = false;
    if (!terminal_handle_in_range(handle)) return 6;
    if ((!input && !resize) || timeout_ms < -1) return 1;
    if (resize && !*initialized) {
        int32_t status = telia_terminal_size(handle, columns, rows);
        if (status) return status;
        *initialized = true;
    }
    uint64_t start = terminal_milliseconds();
    for (;;) {
        if (resize) {
            uint64_t current_columns, current_rows;
            int32_t status = telia_terminal_size(handle, &current_columns, &current_rows);
            if (status) return status;
            *resized = current_columns != *columns || current_rows != *rows;
            *columns = current_columns;
            *rows = current_rows;
        }
#ifdef _WIN32
        if (input) {
            DWORD mode;
            if (GetConsoleMode((HANDLE)handle, &mode)) {
                // ReadFile in canonical mode can still block after a key event.
                if (mode & ENABLE_LINE_INPUT) return 5;
                INPUT_RECORD event;
                DWORD count;
                while (true) {
                    if (!PeekConsoleInputW((HANDLE)handle, &event, 1, &count)) return 6;
                    if (count == 0) break;
                    if (event.EventType == KEY_EVENT && event.Event.KeyEvent.bKeyDown) {
                        *ready = true;
                        break;
                    }
                    if (!ReadConsoleInputW((HANDLE)handle, &event, 1, &count)) return 6;
                }
            } else {
                DWORD type = GetFileType((HANDLE)handle);
                if (type == FILE_TYPE_DISK) *ready = true;
                else if (type == FILE_TYPE_PIPE) {
                    DWORD available;
                    if (!PeekNamedPipe((HANDLE)handle, NULL, 0, NULL, &available, NULL)) {
                        if (GetLastError() == ERROR_BROKEN_PIPE) *ready = true;
                        else return 6;
                    } else *ready = available != 0;
                } else return 5;
            }
        }
        if (*ready || *resized) return 0;
#else
        struct pollfd descriptor = {(int)handle, input ? POLLIN : 0, 0};
        if (!input) descriptor.fd = -1;
        int polled = poll(&descriptor, 1, 0);
        if (polled < 0 && errno != EINTR) return 6;
        if (polled > 0) {
            if (descriptor.revents & (POLLNVAL | POLLERR)) return 6;
            *ready = (descriptor.revents & (POLLIN | POLLHUP)) != 0;
        }
        if (*ready || *resized) return 0;
#endif
        uint64_t elapsed = terminal_milliseconds() - start;
        if (timeout_ms >= 0 && elapsed >= (uint64_t)timeout_ms) {
            *timed_out = true;
            return 0;
        }
        int delay = 16;
        if (timeout_ms >= 0 && (uint64_t)timeout_ms - elapsed < (uint64_t)delay)
            delay = (int)((uint64_t)timeout_ms - elapsed);
#ifdef _WIN32
        Sleep((DWORD)delay);
#else
        int slept = poll(&descriptor, 1, delay);
        if (slept < 0 && errno != EINTR) return 6;
#endif
    }
}
