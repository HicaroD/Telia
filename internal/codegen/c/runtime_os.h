// Private native boundaries for checked filesystem and child-process APIs.
#include <sys/stat.h>
#include <errno.h>
#include <limits.h>
#ifdef _WIN32
#include <io.h>
#include <fcntl.h>
typedef wchar_t _TeliaPathChar;
#else
#include <unistd.h>
#include <fcntl.h>
#include <sys/wait.h>
#include <spawn.h>
extern char **environ;
typedef char _TeliaPathChar;
#endif

static inline _Error _telia_os_error(_Error fallback) {
#ifdef _WIN32
    DWORD code = GetLastError();
    if (code == ERROR_FILE_NOT_FOUND || code == ERROR_PATH_NOT_FOUND)
        return _TELIA_ERROR(NotFound, "path not found");
    if (code == ERROR_ACCESS_DENIED) return _TELIA_ERROR(PermissionDenied, "permission denied");
#else
    if (errno == ENOENT || errno == ENOTDIR) return _TELIA_ERROR(NotFound, "path not found");
    if (errno == EACCES || errno == EPERM) return _TELIA_ERROR(PermissionDenied, "permission denied");
#endif
    return fallback;
}

static inline _Error _telia_native_path(_String path, _TeliaPathChar **result) {
    *result = NULL;
#ifdef _WIN32
    if (!path.len) {
        *result = _TELIA_MALLOC(sizeof(wchar_t));
        if (!*result) return _TELIA_ERROR(AllocationFailed, "allocation failed");
        (*result)[0] = 0;
        return (_Error){0};
    }
    if (path.len > INT_MAX) return _TELIA_ERROR(SizeOverflow, "path too large");
    int count = MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS,
        (const char *)path.data, (int)path.len, NULL, 0);
    if (!count) return _TELIA_ERROR(InvalidUTF8, "invalid UTF-8 path");
    wchar_t *native = _TELIA_MALLOC(((size_t)count + 1) * sizeof(wchar_t));
    if (!native) return _TELIA_ERROR(AllocationFailed, "allocation failed");
    if (!MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS,
            (const char *)path.data, (int)path.len, native, count)) {
        _TELIA_FREE(native);
        return _TELIA_ERROR(InvalidUTF8, "invalid UTF-8 path");
    }
    native[count] = 0;
    *result = native;
    return (_Error){0};
#else
    _String native = {0};
    _Error err = _telia_string_to_c(path, &native);
    *result = (char *)native.data;
    return err;
#endif
}

static inline _Error _telia_path_existing(_String path, _String *result) {
    *result = (_String){0};
    _TeliaPathChar *native = NULL;
    _Error err = _telia_native_path(path, &native);
    if (err.present) return err;
#ifdef _WIN32
    HANDLE file = CreateFileW(native, 0, FILE_SHARE_READ | FILE_SHARE_WRITE | FILE_SHARE_DELETE,
        NULL, OPEN_EXISTING, FILE_FLAG_BACKUP_SEMANTICS, NULL);
    if (file == INVALID_HANDLE_VALUE) {
        err = _telia_os_error(_TELIA_ERROR(PathFailed, "canonicalization failed"));
    } else {
        DWORD count = GetFinalPathNameByHandleW(file, NULL, 0, FILE_NAME_NORMALIZED);
        wchar_t *buffer = count ? _TELIA_MALLOC(((size_t)count + 1) * sizeof(wchar_t)) : NULL;
        if (!count) err = _telia_os_error(_TELIA_ERROR(PathFailed, "canonicalization failed"));
        else if (!buffer) err = _TELIA_ERROR(AllocationFailed, "allocation failed");
        else {
            DWORD written = GetFinalPathNameByHandleW(file, buffer, count + 1, FILE_NAME_NORMALIZED);
            if (!written || written > count) err = _TELIA_ERROR(PathFailed, "canonicalization failed");
            else if (wcsncmp(buffer, L"\\\\?\\UNC\\", 8) == 0) {
                buffer[6] = L'\\';
                err = _telia_text_from_wide(buffer + 6, result);
            } else err = _telia_text_from_wide(wcsncmp(buffer, L"\\\\?\\", 4) == 0 ? buffer + 4 : buffer, result);
        }
        _TELIA_FREE(buffer);
        CloseHandle(file);
    }
#else
    char *resolved = realpath(native, NULL);
    if (!resolved) err = _telia_os_error(_TELIA_ERROR(PathFailed, "canonicalization failed"));
    else {
        err = _telia_string_clone(_telia_string_from_c(resolved), result);
        free(resolved); // realpath allocated this storage, not the runtime allocator.
    }
#endif
    _TELIA_FREE(native);
    return err;
}

static inline _Error _telia_path_missing(_String path, _String *result) {
    _Error err = {0};
    size_t split = path.len;
    while (split && path.data[split - 1] != '/'
#ifdef _WIN32
            && path.data[split - 1] != '\\'
#endif
    ) --split;
    if (split == path.len) return _TELIA_ERROR(InvalidPath, "missing final filename");
    _String parent = split ? (_String){path.data, split} : (_String){(const unsigned char *)".", 1};
    _String resolved = {0};
    err = _telia_path_existing(parent, &resolved);
    if (err.present) return err;
    // Canonical parents must be directories, not merely existing files.
    uint32_t kind = 0;
    _TeliaPathChar *parent_native = NULL;
    err = _telia_native_path(resolved, &parent_native);
    if (!err.present) {
#ifdef _WIN32
        DWORD attrs = GetFileAttributesW(parent_native);
        kind = attrs != INVALID_FILE_ATTRIBUTES && (attrs & FILE_ATTRIBUTE_DIRECTORY) ? 2 : 0;
#else
        struct stat info;
        kind = stat(parent_native, &info) == 0 && S_ISDIR(info.st_mode) ? 2 : 0;
#endif
        if (kind != 2) err = _TELIA_ERROR(NotFound, "parent directory not found");
    }
    _TELIA_FREE(parent_native);
    if (!err.present) {
        _String prefix = {0};
#ifdef _WIN32
        const unsigned char *separator = (const unsigned char *)"\\";
#else
        const unsigned char *separator = (const unsigned char *)"/";
#endif
        bool has_separator = resolved.len && (resolved.data[resolved.len - 1] == '/' || resolved.data[resolved.len - 1] == separator[0]);
        err = _telia_string_concat(resolved, (_String){separator, has_separator ? 0 : 1}, &prefix);
        if (!err.present) err = _telia_string_concat(prefix, (_String){path.data + split, path.len - split}, result);
        _telia_string_free(&prefix);
    }
    _telia_string_free(&resolved);
    return err;
}

static inline _Error _telia_path_canonical(_String path, bool allow_missing, _String *result) {
    _Error err = _telia_path_existing(path, result);
    if (!err.present || !allow_missing || err.kind != &_telia_runtime_kind_NotFound) return err;
    // Never reinterpret a dangling symlink as a new regular file.
    _TeliaPathChar *native = NULL;
    err = _telia_native_path(path, &native);
    if (err.present) return err;
#ifdef _WIN32
    DWORD attrs = GetFileAttributesW(native);
    if (attrs != INVALID_FILE_ATTRIBUTES) {
        _TELIA_FREE(native);
        return _TELIA_ERROR(PathFailed, "dangling reparse point cannot be canonicalized");
    }
    // Expand drive-relative names (C:new) before splitting the missing filename.
    DWORD needed = GetFullPathNameW(native, 0, NULL, NULL);
    wchar_t *buffer = needed ? _TELIA_MALLOC((size_t)needed * sizeof(wchar_t)) : NULL;
    if (!needed) err = _telia_os_error(_TELIA_ERROR(PathFailed, "absolute path lookup failed"));
    else if (!buffer) err = _TELIA_ERROR(AllocationFailed, "allocation failed");
    else {
        DWORD count = GetFullPathNameW(native, needed, buffer, NULL);
        if (!count || count >= needed) err = _TELIA_ERROR(PathFailed, "absolute path lookup failed");
        else {
            _String absolute = {0};
            err = _telia_text_from_wide(buffer, &absolute);
            if (!err.present) err = _telia_path_missing(absolute, result);
            _telia_string_free(&absolute);
        }
    }
    _TELIA_FREE(buffer);
#else
    struct stat info;
    if (lstat(native, &info) == 0) err = _TELIA_ERROR(PathFailed, "dangling symlink cannot be canonicalized");
    else err = _telia_path_missing(path, result);
#endif
    _TELIA_FREE(native);
    return err;
}

static inline _Error _telia_chdir(_String path) {
    _TeliaPathChar *native = NULL;
    _Error err = _telia_native_path(path, &native);
    if (err.present) return err;
#ifdef _WIN32
    if (!SetCurrentDirectoryW(native)) err = _telia_os_error(_TELIA_ERROR(PathFailed, "changing directory failed"));
#else
    if (chdir(native) != 0) err = _telia_os_error(_TELIA_ERROR(PathFailed, "changing directory failed"));
#endif
    _TELIA_FREE(native);
    return err;
}

static inline _Error _telia_fs_metadata(_String path, bool follow, uint32_t *kind,
        uint64_t *volume, uint64_t *identity, uint64_t *size,
        int64_t *seconds, uint32_t *nanos, uint32_t *permissions) {
    _TeliaPathChar *native = NULL;
    _Error err = _telia_native_path(path, &native);
    if (err.present) return err;
#ifdef _WIN32
    HANDLE file = CreateFileW(native, FILE_READ_ATTRIBUTES,
        FILE_SHARE_READ | FILE_SHARE_WRITE | FILE_SHARE_DELETE, NULL, OPEN_EXISTING,
        FILE_FLAG_BACKUP_SEMANTICS | (follow ? 0 : FILE_FLAG_OPEN_REPARSE_POINT), NULL);
    if (file == INVALID_HANDLE_VALUE) {
        err = _telia_os_error(_TELIA_ERROR(MetadataFailed, "metadata lookup failed"));
    } else {
        BY_HANDLE_FILE_INFORMATION info;
        if (!GetFileInformationByHandle(file, &info)) {
            err = _telia_os_error(_TELIA_ERROR(MetadataFailed, "metadata lookup failed"));
        } else {
            *kind = (!follow && (info.dwFileAttributes & FILE_ATTRIBUTE_REPARSE_POINT)) ? 3 :
                (info.dwFileAttributes & FILE_ATTRIBUTE_DIRECTORY) ? 2 :
                GetFileType(file) == FILE_TYPE_DISK ? 1 : 4;
            *volume = info.dwVolumeSerialNumber;
            *identity = ((uint64_t)info.nFileIndexHigh << 32) | info.nFileIndexLow;
            *size = ((uint64_t)info.nFileSizeHigh << 32) | info.nFileSizeLow;
            uint64_t ticks = ((uint64_t)info.ftLastWriteTime.dwHighDateTime << 32) |
                info.ftLastWriteTime.dwLowDateTime;
            *seconds = (int64_t)(ticks / 10000000) - INT64_C(11644473600);
            *nanos = (uint32_t)((ticks % 10000000) * 100);
            *permissions = info.dwFileAttributes;
        }
        CloseHandle(file);
    }
#else
    struct stat info;
    if ((follow ? stat(native, &info) : lstat(native, &info)) != 0) {
        err = _telia_os_error(_TELIA_ERROR(MetadataFailed, "metadata lookup failed"));
    } else {
        *kind = S_ISREG(info.st_mode) ? 1 : S_ISDIR(info.st_mode) ? 2 : S_ISLNK(info.st_mode) ? 3 : 4;
        *volume = (uint64_t)info.st_dev;
        *identity = (uint64_t)info.st_ino;
        *size = info.st_size < 0 ? 0 : (uint64_t)info.st_size;
#ifdef __APPLE__
        *seconds = info.st_mtimespec.tv_sec;
        *nanos = (uint32_t)info.st_mtimespec.tv_nsec;
#else
        *seconds = info.st_mtim.tv_sec;
        *nanos = (uint32_t)info.st_mtim.tv_nsec;
#endif
        *permissions = (uint32_t)(info.st_mode & 07777);
    }
#endif
    _TELIA_FREE(native);
    return err;
}

static inline _Error _telia_file_open(_String path, void **result) {
    *result = NULL;
    _TeliaPathChar *native = NULL;
    _Error err = _telia_native_path(path, &native);
    if (err.present) return err;
#ifdef _WIN32
    FILE *file = _wfopen(native, L"rb");
    if (!file) {
        if (errno == ENOENT || errno == ENOTDIR) err = _TELIA_ERROR(NotFound, "file not found");
        else if (errno == EACCES || errno == EPERM) err = _TELIA_ERROR(PermissionDenied, "permission denied");
        else err = _TELIA_ERROR(FileOpenFailed, "file open failed");
    }
#else
    FILE *file = fopen(native, "rb");
    if (!file) err = _telia_os_error(_TELIA_ERROR(FileOpenFailed, "file open failed"));
    else if (fcntl(fileno(file), F_SETFD, FD_CLOEXEC) != 0) {
        fclose(file);
        file = NULL;
        err = _TELIA_ERROR(FileOpenFailed, "file close-on-exec failed");
    }
#endif
    _TELIA_FREE(native);
    *result = file;
    return err;
}

static inline _Error _telia_file_read(void *state, void *buffer, uint64_t size, uint64_t *count, bool *eof) {
    *count = 0; *eof = false;
    if (!state) return _TELIA_ERROR(InvalidState, "file is closed");
    if (size > SIZE_MAX) return _TELIA_ERROR(SizeOverflow, "read buffer too large");
    if (!size) return (_Error){0};
    *count = fread(buffer, 1, (size_t)size, (FILE *)state);
    *eof = feof((FILE *)state) != 0;
    if (ferror((FILE *)state) || (!*count && !*eof))
        return _TELIA_ERROR(FileReadFailed, "file read failed");
    return (_Error){0};
}

static inline _Error _telia_file_close(void **state) {
    FILE *file = *state;
    *state = NULL;
    if (file && fclose(file) != 0) return _TELIA_ERROR(FileCloseFailed, "file close failed");
    return (_Error){0};
}

static inline _Error _telia_file_preview(_String path, uint64_t limit, _String *result) {
    *result = (_String){0};
    if (limit > SIZE_MAX) return _TELIA_ERROR(SizeOverflow, "preview too large");
    void *file = NULL;
    _Error err = _telia_file_open(path, &file);
    if (err.present) return err;
    unsigned char *data = limit ? _TELIA_MALLOC((size_t)limit) : NULL;
    uint64_t count = 0;
    bool eof = false;
    if (limit && !data) err = _TELIA_ERROR(AllocationFailed, "allocation failed");
    if (!err.present) err = _telia_file_read(file, data, limit, &count, &eof);
    _Error close_err = _telia_file_close(&file);
    if (!err.present) err = close_err;
    if (err.present) _TELIA_FREE(data);
    else *result = (_String){data, (size_t)count};
    return err;
}

#ifndef _TELIA_FILE_SYNC
#ifdef _WIN32
#define _TELIA_FILE_SYNC _commit
#else
#define _TELIA_FILE_SYNC fsync
#endif
#endif
#ifndef _TELIA_FILE_RENAME
#define _TELIA_FILE_RENAME rename
#endif

static inline _Error _telia_file_remove(_String path) {
    _TeliaPathChar *native = NULL;
    _Error err = _telia_native_path(path, &native);
    if (err.present) return err;
#ifdef _WIN32
    if (!DeleteFileW(native)) err = _telia_os_error(_TELIA_ERROR(FileRemoveFailed, "file removal failed"));
#else
    if (unlink(native) != 0) err = _telia_os_error(_TELIA_ERROR(FileRemoveFailed, "file removal failed"));
#endif
    _TELIA_FREE(native);
    return err;
}

static inline _Error _telia_file_temp(_String destination, void **state, _String *name) {
    *state = NULL; *name = (_String){0};
    _TeliaPathChar *native = NULL;
    _Error err = _telia_native_path(destination, &native);
    if (err.present) return err;
    int fd = -1;
#ifdef _WIN32
    size_t len = wcslen(native);
    wchar_t *pattern = _TELIA_MALLOC((len + 80) * sizeof(wchar_t));
    if (!pattern) { _TELIA_FREE(native); return _TELIA_ERROR(AllocationFailed, "allocation failed"); }
    static LONG sequence;
    for (int attempt = 0; attempt < 128; ++attempt) {
        _snwprintf(pattern, len + 80, L"%ls.telia-%lu-%lu-%lu", native,
            (unsigned long)GetCurrentProcessId(), (unsigned long)GetTickCount(),
            (unsigned long)InterlockedIncrement(&sequence));
        fd = _wopen(pattern, _O_CREAT | _O_EXCL | _O_RDWR | _O_BINARY | _O_NOINHERIT, _S_IREAD | _S_IWRITE);
        if (fd >= 0 || errno != EEXIST) break;
    }
    if (fd < 0) {
        if (errno == ENOENT || errno == ENOTDIR) err = _TELIA_ERROR(NotFound, "temporary directory not found");
        else if (errno == EACCES || errno == EPERM) err = _TELIA_ERROR(PermissionDenied, "permission denied");
        else err = _TELIA_ERROR(FileOpenFailed, "temporary file creation failed");
    }
    else {
        FILE *file = _wfdopen(fd, L"w+b");
        if (!file) { _close(fd); err = _TELIA_ERROR(FileOpenFailed, "temporary stream creation failed"); }
        else {
            err = _telia_text_from_wide(pattern, name);
            if (err.present) fclose(file);
            else *state = file;
        }
        if (err.present) DeleteFileW(pattern);
    }
#else
    size_t len = strlen(native);
    if (len > SIZE_MAX - 14) { _TELIA_FREE(native); return _TELIA_ERROR(SizeOverflow, "path too large"); }
    char *pattern = _TELIA_MALLOC(len + 14);
    if (!pattern) { _TELIA_FREE(native); return _TELIA_ERROR(AllocationFailed, "allocation failed"); }
    memcpy(pattern, native, len);
    memcpy(pattern + len, ".telia-XXXXXX", 14);
    fd = mkstemp(pattern);
    if (fd < 0) err = _telia_os_error(_TELIA_ERROR(FileOpenFailed, "temporary file creation failed"));
    else {
        FILE *file = NULL;
        if (fcntl(fd, F_SETFD, FD_CLOEXEC) == 0) file = fdopen(fd, "w+b");
        if (!file) { close(fd); err = _TELIA_ERROR(FileOpenFailed, "temporary stream creation failed"); }
        else {
            err = _telia_string_clone(_telia_string_from_c(pattern), name);
            if (err.present) fclose(file);
            else *state = file;
        }
        if (err.present) unlink(pattern);
    }
#endif
    _TELIA_FREE(pattern);
    _TELIA_FREE(native);
    return err;
}

static inline _Error _telia_file_write(void *state, _String data) {
    if (!state) return _TELIA_ERROR(InvalidState, "file is closed");
    size_t offset = 0;
    while (offset < data.len) {
        size_t count = fwrite(data.data + offset, 1, data.len - offset, state);
        offset += count;
        if (!count || ferror(state)) return _TELIA_ERROR(FileWriteFailed, "file write failed");
    }
    return (_Error){0};
}

// Allocation-free cleanup; bounded copies also keep this safe for string views.
static inline _Error _telia_file_cleanup_temp(_String path) {
#ifdef _WIN32
    wchar_t native[32768];
    if (path.len > INT_MAX) return _TELIA_ERROR(SizeOverflow, "temporary path too large");
    int count = MultiByteToWideChar(CP_UTF8, MB_ERR_INVALID_CHARS,
        (const char *)path.data, (int)path.len, native, 32767);
    if (!count) return _TELIA_ERROR(FileRemoveFailed, "temporary path conversion failed");
    native[count] = 0;
    if (!DeleteFileW(native)) return _telia_os_error(_TELIA_ERROR(FileRemoveFailed, "temporary cleanup failed"));
#else
    char native[PATH_MAX];
    if (path.len >= sizeof(native)) return _TELIA_ERROR(SizeOverflow, "temporary path too large");
    if (path.len) memcpy(native, path.data, path.len);
    native[path.len] = 0;
    if (unlink(native) != 0) return _telia_os_error(_TELIA_ERROR(FileRemoveFailed, "temporary cleanup failed"));
#endif
    return (_Error){0};
}

static inline _Error _telia_file_sync(void *state) {
    if (!state) return _TELIA_ERROR(InvalidState, "file is closed");
    if (fflush(state) != 0) return _TELIA_ERROR(FileSyncFailed, "file flush failed");
#ifdef _WIN32
    int fd = _fileno(state);
#else
    int fd = fileno(state);
#endif
    int status;
    do { status = _TELIA_FILE_SYNC(fd); } while (status < 0 && errno == EINTR);
    if (status != 0) return _TELIA_ERROR(FileSyncFailed, "file synchronization failed");
    return (_Error){0};
}

static inline _Error _telia_file_permissions(_String source, void *state) {
    if (!state) return _TELIA_ERROR(InvalidState, "file is closed");
    _TeliaPathChar *native = NULL;
    _Error err = _telia_native_path(source, &native);
    if (err.present) return err;
#ifdef _WIN32
    // ReplaceFileW preserves the existing target's DACL at replacement time.
    // This primitive alone cannot transfer a Windows security descriptor.
    (void)state;
    err = _TELIA_ERROR(Unsupported, "standalone Windows permission copying is unsupported");
#else
    struct stat info;
    if (stat(native, &info) != 0) err = _telia_os_error(_TELIA_ERROR(MetadataFailed, "permission lookup failed"));
    else if (fchmod(fileno(state), info.st_mode & 07777) != 0)
        err = _telia_os_error(_TELIA_ERROR(PermissionDenied, "permission preservation failed"));
#endif
    _TELIA_FREE(native);
    return err;
}

static inline _Error _telia_file_rename(_String source, _String destination) {
    _TeliaPathChar *from = NULL, *to = NULL;
    _Error err = _telia_native_path(source, &from);
    if (!err.present) err = _telia_native_path(destination, &to);
    if (!err.present) {
#ifdef _WIN32
        DWORD attributes = GetFileAttributesW(to);
        BOOL ok = attributes == INVALID_FILE_ATTRIBUTES && GetLastError() == ERROR_FILE_NOT_FOUND ?
            MoveFileExW(from, to, MOVEFILE_WRITE_THROUGH) : ReplaceFileW(to, from, NULL, 0, NULL, NULL);
        if (!ok) err = _telia_os_error(_TELIA_ERROR(FileRenameFailed, "file replacement failed"));
#else
        if (_TELIA_FILE_RENAME(from, to) != 0) err = _telia_os_error(_TELIA_ERROR(FileRenameFailed, "file replacement failed"));
#endif
    }
    _TELIA_FREE(from); _TELIA_FREE(to);
    return err;
}

static inline _Error _telia_parent_sync(_String path) {
#ifdef _WIN32
    (void)path;
    return _TELIA_ERROR(Unsupported, "directory synchronization is unsupported on Windows");
#else
    size_t split = path.len;
    while (split && path.data[split - 1] != '/') --split;
    _String parent = split ? (_String){path.data, split} : (_String){(const unsigned char *)".", 1};
    char *native = NULL;
    _Error err = _telia_native_path(parent, &native);
    if (err.present) return err;
    int fd = open(native, O_RDONLY | O_DIRECTORY | O_CLOEXEC);
    if (fd < 0) err = _telia_os_error(_TELIA_ERROR(FileSyncFailed, "parent directory open failed"));
    else {
        int status;
        do { status = _TELIA_FILE_SYNC(fd); } while (status < 0 && errno == EINTR);
        if (status != 0) err = _TELIA_ERROR(FileSyncFailed, "directory synchronization failed");
        if (close(fd) != 0 && !err.present) err = _TELIA_ERROR(FileCloseFailed, "directory close failed");
    }
    _TELIA_FREE(native);
    return err;
#endif
}
