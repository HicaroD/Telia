#include "allocation_control.h"
#include <stdio.h>
#ifdef _WIN32
#include <io.h>
#include <fcntl.h>
#else
#include <unistd.h>
#include <signal.h>
#include <pthread.h>
#include <fcntl.h>
#endif

static int io_mode, io_calls, io_failed, io_flushes, io_modes;
static void _test_io_mode(int32_t mode) {
    io_mode = mode;
    io_calls = io_failed = io_flushes = io_modes = 0;
}
static int32_t _test_io_flushes(void) { return io_flushes; }
static int32_t _test_io_streams_open(void) {
    // A successful descriptor lookup demonstrates none of the streams was closed.
#ifdef _WIN32
    return _get_osfhandle(_fileno(stdin)) != -1 &&
        _get_osfhandle(_fileno(stdout)) != -1 && _get_osfhandle(_fileno(stderr)) != -1;
#else
    return fcntl(fileno(stdin), F_GETFD) != -1 &&
        fcntl(fileno(stdout), F_GETFD) != -1 && fcntl(fileno(stderr), F_GETFD) != -1;
#endif
}
static size_t _test_io_read(void *data, size_t size, size_t count, FILE *stream) {
    ++io_calls;
    if (io_mode == 2 && io_calls == 2) { io_failed = 1; return 0; }
    if (io_mode == 3) return 0;
    if (io_mode == 1 || io_mode == 2) { if (count > 7) count = 7; }
    size_t read = fread(data, size, count, stream);
    if (io_mode == 11) io_failed = 1;
    return read;
}
static size_t _test_io_write(const void *data, size_t size, size_t count, FILE *stream) {
    ++io_calls;
    if (io_mode == 5) return 0;
    if (io_mode == 1 || io_mode == 4) { if (count > 7) count = 7; }
    size_t written = fwrite(data, size, count, stream);
    if (io_mode == 4) io_failed = 1;
    return written;
}
static int _test_io_error(FILE *stream) { return io_failed || ferror(stream); }
static int _test_io_flush(FILE *stream) {
    ++io_flushes;
    int result = fflush(stream);
    // Windows pre-mode flush must succeed to exercise the final flush.
    if ((io_mode == 6 || io_mode == 4) && io_calls) return EOF;
    if (io_mode == 10) return EOF;
    return result;
}
#ifdef _WIN32
static int _test_io_setmode(int fd, int mode) {
    ++io_modes;
    if (io_mode == 8 || (io_mode == 9 && io_modes == 2)) return -1;
    return _setmode(fd, mode);
}
#define _TELIA_IO_SETMODE _test_io_setmode
#endif
#define _TELIA_IO_READ _test_io_read
#define _TELIA_IO_WRITE _test_io_write
#define _TELIA_IO_ERROR _test_io_error
#define _TELIA_IO_FLUSH _test_io_flush

#ifndef _WIN32
static sigset_t io_original_mask, io_before_pending;
static int io_pending_expected;
static int32_t _test_io_broken_pipe(void) {
    int descriptors[2];
    if (pipe(descriptors) != 0) return 0;
    close(descriptors[0]);
    if (dup2(descriptors[1], fileno(stdout)) == -1) { close(descriptors[1]); return 0; }
    close(descriptors[1]);
    signal(SIGPIPE, SIG_DFL);
    return pthread_sigmask(SIG_SETMASK, NULL, &io_original_mask) == 0;
}
static int32_t _test_io_signal_unchanged(void) {
    sigset_t current, pending;
    struct sigaction action;
    if (pthread_sigmask(SIG_SETMASK, NULL, &current) != 0 ||
        sigpending(&pending) != 0 || sigaction(SIGPIPE, NULL, &action) != 0) return 0;
    int ok = sigismember(&current, SIGPIPE) == sigismember(&io_original_mask, SIGPIPE) &&
        sigismember(&pending, SIGPIPE) == io_pending_expected && action.sa_handler == SIG_DFL;
    if (io_pending_expected) {
        sigset_t blocked;
        int number;
        sigemptyset(&blocked);
        sigaddset(&blocked, SIGPIPE);
        if (sigismember(&pending, SIGPIPE)) sigwait(&blocked, &number);
        pthread_sigmask(SIG_SETMASK, &io_before_pending, NULL);
    }
    return ok;
}
static void _test_io_pending_signal(void) {
    sigset_t blocked;
    sigemptyset(&blocked);
    sigaddset(&blocked, SIGPIPE);
    pthread_sigmask(SIG_BLOCK, &blocked, &io_before_pending);
    raise(SIGPIPE);
    pthread_sigmask(SIG_SETMASK, NULL, &io_original_mask);
    io_pending_expected = 1;
}
#endif

static void _test_io_native_text(void) {
    fputs("native\n", stdout);
    fputs("native\n", stderr);
}
