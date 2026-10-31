#include <stdio.h>
#include <wchar.h>
#include "allocation_control.h"
#ifdef _WIN32
#include <fcntl.h>
#include <io.h>
#else
#include <unistd.h>
#endif

static int fs_live = 0, fs_opens = 0, fs_mode = 0, fs_reads = 0, fs_writes = 0;
static void _test_fs_mode(int32_t mode) {
    fs_mode = mode;
    fs_reads = fs_writes = 0;
}
static int32_t _test_fs_handles(void) { return fs_live; }
static int32_t _test_fs_opens(void) { return fs_opens; }
static FILE *test_fs_pipe(void) {
    int ends[2];
    const unsigned char bytes[] = {'A', 0, 255, '\r', '\n'};
#ifdef _WIN32
    if (_pipe(ends, 256, _O_BINARY) != 0) abort();
    if (_write(ends[1], bytes, sizeof(bytes)) != sizeof(bytes)) abort();
    _close(ends[1]);
    FILE *file = _fdopen(ends[0], "rb");
#else
    if (pipe(ends) != 0) abort();
    if (write(ends[1], bytes, sizeof(bytes)) != sizeof(bytes)) abort();
    close(ends[1]);
    FILE *file = fdopen(ends[0], "rb");
#endif
    if (!file) abort();
    ++fs_live;
    return file;
}
#ifdef _WIN32
static FILE *test_fs_open(const wchar_t *path, const wchar_t *mode) {
    ++fs_opens;
    if (fs_mode == 7) return NULL;
    if (fs_mode == 8) return test_fs_pipe();
    FILE *file = _wfopen(path, mode);
#else
static FILE *test_fs_open(const char *path, const char *mode) {
    ++fs_opens;
    if (fs_mode == 7) return NULL;
    if (fs_mode == 8) return test_fs_pipe();
    FILE *file = fopen(path, mode);
#endif
    if (file) ++fs_live;
    return file;
}
static size_t test_fs_read(void *data, size_t size, size_t count, FILE *file) {
    ++fs_reads;
    if ((fs_mode == 2 || fs_mode == 5) && fs_reads > 1) return 0;
    if (fs_mode == 1 && count > 7) count = 7;
    return fread(data, size, count, file);
}
static size_t test_fs_write(const void *data, size_t size, size_t count, FILE *file) {
    ++fs_writes;
    if ((fs_mode == 3 || fs_mode == 6) && fs_writes > 1) return 0;
    if ((fs_mode == 1 || fs_mode == 3 || fs_mode == 6) && count > 7) count = 7;
    return fwrite(data, size, count, file);
}
static int test_fs_error(FILE *file) {
    if ((fs_mode == 2 && fs_reads > 1) || (fs_mode == 3 && fs_writes > 1)) return 1;
    return ferror(file);
}
static int test_fs_close(FILE *file) {
    --fs_live;
    int result = fclose(file);
    // Also fail cleanup after operation failures to check precedence.
    return fs_mode >= 2 && fs_mode <= 4 ? EOF : result;
}
__attribute__((destructor)) static void test_fs_no_handles(void) {
    if (fs_live != 0) abort();
}
#define _TELIA_FS_OPEN test_fs_open
#define _TELIA_FS_READ test_fs_read
#define _TELIA_FS_WRITE test_fs_write
#define _TELIA_FS_ERROR test_fs_error
#define _TELIA_FS_CLOSE test_fs_close
