#include <errno.h>
#ifdef _WIN32
#include <io.h>
#else
#include <unistd.h>
#endif

static int _test_save_failure;
static int _test_sync_calls;
static void _test_save_mode(int mode) {
    _test_save_failure = mode;
    _test_sync_calls = 0;
}

static int _test_save_sync(int fd) {
    ++_test_sync_calls;
    if (_test_save_failure == 1 || (_test_save_failure == 3 && _test_sync_calls == 2)) {
        errno = EIO;
        return -1;
    }
#ifdef _WIN32
    return _commit(fd);
#else
    return fsync(fd);
#endif
}

#ifndef _WIN32
#include <stdio.h>
static int _test_save_rename(const char *source, const char *destination) {
    if (_test_save_failure == 2) { errno = EIO; return -1; }
    return rename(source, destination);
}
#define _TELIA_FILE_RENAME _test_save_rename
#endif
#define _TELIA_FILE_SYNC _test_save_sync
