// Owning native child state; all public failure results are checked errors.
typedef struct {
#ifdef _WIN32
    HANDLE process;
    HANDLE pipes[3];
#else
    pid_t pid;
    int pipes[3];
#endif
    bool done;
    bool eof[3];
    uint32_t code, signal_number;
} _TeliaChild;

static inline _Error _telia_child_close_pipes(_TeliaChild *child) {
    _Error err = {0};
    for (int i = 0; i < 3; ++i) {
#ifdef _WIN32
        if (child->pipes[i] && !CloseHandle(child->pipes[i]) && !err.present)
            err = _TELIA_ERROR(FileCloseFailed, "pipe close failed");
        child->pipes[i] = NULL;
#else
        if (child->pipes[i] >= 0 && close(child->pipes[i]) != 0 && !err.present)
            err = _TELIA_ERROR(FileCloseFailed, "pipe close failed");
        child->pipes[i] = -1;
#endif
    }
    return err;
}

#ifndef _WIN32
static inline int _telia_child_pipe(int pair[2], bool input) {
    int raw[2];
    if (pipe(raw) != 0) return -1;
    // Never let a pipe collide with stdin/out/err, even when the caller closed them.
    for (int i = 0; i < 2; ++i) {
        pair[i] = fcntl(raw[i], F_DUPFD_CLOEXEC, 3);
        if (pair[i] < 0) {
            if (i) close(pair[0]);
            close(raw[0]); close(raw[1]);
            pair[0] = pair[1] = -1;
            return -1;
        }
    }
    close(raw[0]); close(raw[1]);
    int parent = pair[input ? 1 : 0];
    if (fcntl(parent, F_SETFL, O_NONBLOCK) < 0) {
        close(pair[0]); close(pair[1]);
        pair[0] = pair[1] = -1;
        return -1;
    }
    return 0;
}
#else
static inline _Error _telia_child_command(wchar_t **argv, size_t count, wchar_t **result) {
    *result = NULL;
    size_t capacity = 1;
    for (size_t i = 0; i < count; ++i) {
        size_t len = wcslen(argv[i]);
        if (len > 16380 || capacity + len * 2 + 3 > 32767)
            return _TELIA_ERROR(SizeOverflow, "Windows command line too large");
        capacity += len * 2 + 3;
    }
    wchar_t *command = _TELIA_MALLOC(capacity * sizeof(wchar_t));
    if (!command) return _TELIA_ERROR(AllocationFailed, "allocation failed");
    size_t at = 0;
    for (size_t i = 0; i < count; ++i) {
        if (i) command[at++] = L' ';
        command[at++] = L'"';
        size_t slashes = 0;
        for (const wchar_t *p = argv[i]; ; ++p) {
            if (*p == L'\\') { ++slashes; continue; }
            size_t copies = (*p == L'"' || !*p) ? slashes * 2 : slashes;
            while (copies--) command[at++] = L'\\';
            slashes = 0;
            if (!*p) break;
            if (*p == L'"') command[at++] = L'\\';
            command[at++] = *p;
        }
        command[at++] = L'"';
    }
    command[at] = 0;
    *result = command;
    return (_Error){0};
}

static inline bool _telia_child_windows_pipe(int index, HANDLE *parent, HANDLE *child) {
    SECURITY_ATTRIBUTES security = {sizeof(security), NULL, TRUE};
    if (index != 0) {
        if (!CreatePipe(parent, child, &security, 0)) return false;
        if (!SetHandleInformation(*parent, HANDLE_FLAG_INHERIT, 0)) {
            CloseHandle(*parent); CloseHandle(*child); *parent = *child = NULL;
            return false;
        }
        return true;
    }
    // PIPE_NOWAIT makes stdin writes partial/nonblocking too, avoiding pipe deadlocks.
    wchar_t name[128];
    static LONG sequence;
    _snwprintf(name, 128, L"\\\\.\\pipe\\telia-%lu-%lu-%lu", (unsigned long)GetCurrentProcessId(),
        (unsigned long)GetTickCount(), (unsigned long)InterlockedIncrement(&sequence));
    *parent = CreateNamedPipeW(name, PIPE_ACCESS_OUTBOUND | FILE_FLAG_FIRST_PIPE_INSTANCE,
        PIPE_TYPE_BYTE | PIPE_READMODE_BYTE | PIPE_NOWAIT, 1, 65536, 65536, 0, NULL);
    if (*parent == INVALID_HANDLE_VALUE) { *parent = NULL; return false; }
    *child = CreateFileW(name, GENERIC_READ, 0, &security, OPEN_EXISTING, 0, NULL);
    if (*child == INVALID_HANDLE_VALUE) {
        CloseHandle(*parent); *parent = *child = NULL;
        return false;
    }
    if (!ConnectNamedPipe(*parent, NULL) && GetLastError() != ERROR_PIPE_CONNECTED) {
        CloseHandle(*parent); CloseHandle(*child); *parent = *child = NULL;
        return false;
    }
    return true;
}
#endif

static inline _Error _telia_child_spawn(_String executable, void *arguments,
        uint64_t count, bool input, bool output, bool stderr_pipe, void **result) {
    *result = NULL;
    if (count > SIZE_MAX / sizeof(_TeliaPathChar *) - 2)
        return _TELIA_ERROR(SizeOverflow, "too many process arguments");
    _TeliaChild *child = _TELIA_MALLOC(sizeof(*child));
    if (!child) return _TELIA_ERROR(AllocationFailed, "allocation failed");
    *child = (_TeliaChild){0};
#ifndef _WIN32
    child->pipes[0] = child->pipes[1] = child->pipes[2] = -1;
#endif
    _TeliaPathChar **argv = _TELIA_MALLOC(((size_t)count + 2) * sizeof(*argv));
    if (!argv) { _TELIA_FREE(child); return _TELIA_ERROR(AllocationFailed, "allocation failed"); }
    memset(argv, 0, ((size_t)count + 2) * sizeof(*argv));
    _Error err = _telia_native_path(executable, &argv[0]);
    for (size_t i = 0; !err.present && i < count; ++i)
        err = _telia_native_path(((_String *)arguments)[i], &argv[i + 1]);
    bool requested[3] = {input, output, stderr_pipe};
#ifdef _WIN32
    HANDLE inherited[3] = {0};
    wchar_t *command = NULL;
    LPPROC_THREAD_ATTRIBUTE_LIST attributes = NULL;
    bool attributes_ready = false;
    if (!err.present) err = _telia_child_command(argv, (size_t)count + 1, &command);
    for (int i = 0; !err.present && i < 3; ++i) {
        if (requested[i]) {
            if (!_telia_child_windows_pipe(i, &child->pipes[i], &inherited[i]))
                err = _TELIA_ERROR(ProcessSpawnFailed, "process pipe creation failed");
        } else {
            HANDLE source = GetStdHandle(i == 0 ? STD_INPUT_HANDLE : i == 1 ? STD_OUTPUT_HANDLE : STD_ERROR_HANDLE);
            if (!source || source == INVALID_HANDLE_VALUE) {
                SECURITY_ATTRIBUTES security = {sizeof(security), NULL, TRUE};
                inherited[i] = CreateFileW(L"NUL", i == 0 ? GENERIC_READ : GENERIC_WRITE,
                    FILE_SHARE_READ | FILE_SHARE_WRITE, &security, OPEN_EXISTING, 0, NULL);
                if (inherited[i] == INVALID_HANDLE_VALUE) { inherited[i] = NULL; err = _TELIA_ERROR(ProcessSpawnFailed, "standard handle creation failed"); }
            } else if (!DuplicateHandle(GetCurrentProcess(), source, GetCurrentProcess(), &inherited[i], 0, TRUE, DUPLICATE_SAME_ACCESS))
                err = _TELIA_ERROR(ProcessSpawnFailed, "standard handle duplication failed");
        }
    }
    if (!err.present) {
        SIZE_T size = 0;
        InitializeProcThreadAttributeList(NULL, 1, 0, &size);
        attributes = _TELIA_MALLOC(size);
        if (!attributes) err = _TELIA_ERROR(AllocationFailed, "allocation failed");
        else if (!InitializeProcThreadAttributeList(attributes, 1, 0, &size))
            err = _TELIA_ERROR(ProcessSpawnFailed, "process attributes failed");
        else {
            attributes_ready = true;
            if (!UpdateProcThreadAttribute(attributes, 0, PROC_THREAD_ATTRIBUTE_HANDLE_LIST,
                    inherited, sizeof(inherited), NULL, NULL))
                err = _TELIA_ERROR(ProcessSpawnFailed, "process handle inheritance failed");
        }
    }
    if (!err.present) {
        STARTUPINFOEXW startup = {0};
        startup.StartupInfo.cb = sizeof(startup);
        startup.StartupInfo.dwFlags = STARTF_USESTDHANDLES;
        startup.StartupInfo.hStdInput = inherited[0];
        startup.StartupInfo.hStdOutput = inherited[1];
        startup.StartupInfo.hStdError = inherited[2];
        startup.lpAttributeList = attributes;
        PROCESS_INFORMATION info = {0};
        if (!CreateProcessW(NULL, command, NULL, NULL, TRUE, EXTENDED_STARTUPINFO_PRESENT,
                NULL, NULL, &startup.StartupInfo, &info))
            err = _TELIA_ERROR(ProcessSpawnFailed, "process spawn failed");
        else { child->process = info.hProcess; CloseHandle(info.hThread); }
    }
    if (attributes_ready) DeleteProcThreadAttributeList(attributes);
    _TELIA_FREE(attributes); _TELIA_FREE(command);
    for (int i = 0; i < 3; ++i) if (inherited[i]) CloseHandle(inherited[i]);
#else
    int pairs[3][2] = {{-1, -1}, {-1, -1}, {-1, -1}};
    posix_spawn_file_actions_t actions;
    bool actions_ready = false;
    if (!err.present) {
        if (posix_spawn_file_actions_init(&actions) != 0) err = _TELIA_ERROR(ProcessSpawnFailed, "process actions failed");
        else actions_ready = true;
    }
    for (int i = 0; !err.present && i < 3; ++i) {
        if (!requested[i]) continue;
        if (_telia_child_pipe(pairs[i], i == 0) != 0) { err = _TELIA_ERROR(ProcessSpawnFailed, "process pipe creation failed"); break; }
        int native = pairs[i][i == 0 ? 0 : 1];
        if (posix_spawn_file_actions_adddup2(&actions, native, i) != 0 ||
                posix_spawn_file_actions_addclose(&actions, pairs[i][0]) != 0 ||
                posix_spawn_file_actions_addclose(&actions, pairs[i][1]) != 0)
            err = _TELIA_ERROR(ProcessSpawnFailed, "process pipe actions failed");
    }
    if (!err.present && posix_spawnp(&child->pid, argv[0], &actions, NULL, argv, environ) != 0)
        err = _TELIA_ERROR(ProcessSpawnFailed, "process spawn failed");
    if (actions_ready) posix_spawn_file_actions_destroy(&actions);
    for (int i = 0; i < 3; ++i) {
        int parent = i == 0 ? 1 : 0;
        if (pairs[i][1 - parent] >= 0) close(pairs[i][1 - parent]);
        child->pipes[i] = pairs[i][parent];
    }
#endif
    for (size_t i = 0; i < (size_t)count + 1; ++i) _TELIA_FREE(argv[i]);
    _TELIA_FREE(argv);
    if (err.present) { _telia_child_close_pipes(child); _TELIA_FREE(child); }
    else *result = child;
    return err;
}

static inline _Error _telia_child_wait(void *state, bool blocking, bool *done, uint32_t *code, uint32_t *signal_number) {
    *done = false; *code = 0; *signal_number = 0;
    if (!state) return _TELIA_ERROR(InvalidState, "child is closed");
    _TeliaChild *child = state;
    if (!child->done) {
#ifdef _WIN32
        DWORD status = WaitForSingleObject(child->process, blocking ? INFINITE : 0);
        if (status == WAIT_TIMEOUT) return (_Error){0};
        DWORD exit_code;
        if (status != WAIT_OBJECT_0 || !GetExitCodeProcess(child->process, &exit_code))
            return _TELIA_ERROR(ProcessWaitFailed, "process wait failed");
        child->code = exit_code;
#else
        int status;
        pid_t pid;
        do { pid = waitpid(child->pid, &status, blocking ? 0 : WNOHANG); } while (pid < 0 && errno == EINTR);
        if (pid == 0) return (_Error){0};
        if (pid < 0) return _TELIA_ERROR(ProcessWaitFailed, "process wait failed");
        if (WIFEXITED(status)) child->code = (uint32_t)WEXITSTATUS(status);
        else if (WIFSIGNALED(status)) child->signal_number = (uint32_t)WTERMSIG(status);
        else return _TELIA_ERROR(ProcessWaitFailed, "unexpected process status");
#endif
        child->done = true;
    }
    *done = true; *code = child->code; *signal_number = child->signal_number;
    return (_Error){0};
}

static inline _Error _telia_child_terminate(void *state, bool force) {
    if (!state) return _TELIA_ERROR(InvalidState, "child is closed");
    bool done; uint32_t code, signal_number;
    _Error err = _telia_child_wait(state, false, &done, &code, &signal_number);
    if (err.present || done) return err;
    _TeliaChild *child = state;
#ifdef _WIN32
    (void)force;
    if (!TerminateProcess(child->process, 1) && WaitForSingleObject(child->process, 0) != WAIT_OBJECT_0)
        return _TELIA_ERROR(ProcessTerminateFailed, "process termination failed");
#else
    if (kill(child->pid, force ? SIGKILL : SIGTERM) != 0 && errno != ESRCH)
        return _TELIA_ERROR(ProcessTerminateFailed, "process termination failed");
#endif
    return (_Error){0};
}

static inline _Error _telia_child_dispose(void **state) {
    if (!*state) return (_Error){0};
    _TeliaChild *child = *state;
    _Error err = _telia_child_close_pipes(child);
    _Error terminate_err = _telia_child_terminate(child, true);
    if (terminate_err.present) return err.present ? err : terminate_err;
    bool done; uint32_t code, signal_number;
    _Error wait_err = _telia_child_wait(child, true, &done, &code, &signal_number);
    if (!err.present) err = wait_err;
    // Retain ownership if reaping failed, so callers may retry cleanup.
    if (!child->done) return err;
#ifdef _WIN32
    if (!CloseHandle(child->process) && !err.present)
        err = _TELIA_ERROR(FileCloseFailed, "process handle close failed");
#endif
    _TELIA_FREE(child); *state = NULL;
    return err;
}

static inline _Error _telia_child_close_pipe(void *state, uint32_t stream) {
    if (!state || stream > 2) return _TELIA_ERROR(InvalidState, "invalid child pipe");
    _TeliaChild *child = state;
    int status = 0;
#ifdef _WIN32
    if (child->pipes[stream]) status = CloseHandle(child->pipes[stream]) ? 0 : -1;
    child->pipes[stream] = NULL;
#else
    if (child->pipes[stream] >= 0) status = close(child->pipes[stream]);
    child->pipes[stream] = -1;
#endif
    if (status != 0) return _TELIA_ERROR(FileCloseFailed, "pipe close failed");
    return (_Error){0};
}

static inline _Error _telia_child_read(void *state, uint32_t stream, void *buffer, uint64_t size,
        uint64_t *count, bool *eof, bool *would_block) {
    *count = 0; *eof = false; *would_block = false;
    if (!state || stream < 1 || stream > 2) return _TELIA_ERROR(InvalidState, "invalid child pipe");
    if (size > SIZE_MAX) return _TELIA_ERROR(SizeOverflow, "pipe read buffer too large");
    _TeliaChild *child = state;
    if (child->eof[stream]) { *eof = true; return (_Error){0}; }
#ifdef _WIN32
    HANDLE pipe = child->pipes[stream];
    if (!pipe) return _TELIA_ERROR(InvalidState, "output pipe is closed");
    if (!size) return (_Error){0};
    DWORD available = 0, read_count = 0;
    if (!PeekNamedPipe(pipe, NULL, 0, NULL, &available, NULL)) {
        if (GetLastError() != ERROR_BROKEN_PIPE) return _TELIA_ERROR(StreamReadFailed, "pipe probe failed");
        *eof = true;
    } else if (!available) *would_block = true;
    else {
        DWORD request = size < available ? (DWORD)size : available;
        if (!ReadFile(pipe, buffer, request, &read_count, NULL)) {
            if (GetLastError() != ERROR_BROKEN_PIPE) return _TELIA_ERROR(StreamReadFailed, "pipe read failed");
            *eof = true;
        } else *count = read_count;
    }
#else
    int pipe = child->pipes[stream];
    if (pipe < 0) return _TELIA_ERROR(InvalidState, "output pipe is closed");
    if (!size) return (_Error){0};
    ssize_t read_count;
    do { read_count = read(pipe, buffer, size > 65536 ? 65536 : (size_t)size); }
        while (read_count < 0 && errno == EINTR);
    if (read_count > 0) *count = (uint64_t)read_count;
    else if (!read_count) *eof = true;
    else if (errno == EAGAIN || errno == EWOULDBLOCK) *would_block = true;
    else return _TELIA_ERROR(StreamReadFailed, "pipe read failed");
#endif
    if (*eof) {
        child->eof[stream] = true;
        return _telia_child_close_pipe(child, stream);
    }
    return (_Error){0};
}

static inline _Error _telia_child_write(void *state, _String data, uint64_t *count, bool *would_block) {
    *count = 0; *would_block = false;
    if (!state) return _TELIA_ERROR(InvalidState, "child is closed");
    _TeliaChild *child = state;
#ifdef _WIN32
    if (!child->pipes[0]) return _TELIA_ERROR(InvalidState, "stdin pipe is closed");
    if (!data.len) return (_Error){0};
    DWORD written = 0;
    if (!WriteFile(child->pipes[0], data.data, data.len > 65536 ? 65536 : (DWORD)data.len, &written, NULL)) {
        return _TELIA_ERROR(StreamWriteFailed, "stdin pipe write failed");
    }
    *count = written;
    *would_block = written == 0;
    return (_Error){0};
#else
    if (child->pipes[0] < 0) return _TELIA_ERROR(InvalidState, "stdin pipe is closed");
    if (!data.len) return (_Error){0};
    // A broken child pipe must be an error, not a fatal signal to the caller.
    sigset_t blocked, previous, pending;
    sigemptyset(&blocked); sigaddset(&blocked, SIGPIPE);
    if (pthread_sigmask(SIG_BLOCK, &blocked, &previous) != 0)
        return _TELIA_ERROR(StreamWriteFailed, "cannot block SIGPIPE");
    if (sigpending(&pending) != 0) {
        pthread_sigmask(SIG_SETMASK, &previous, NULL);
        return _TELIA_ERROR(StreamWriteFailed, "cannot inspect pending signals");
    }
    bool already_pending = sigismember(&pending, SIGPIPE) == 1;
    ssize_t written;
    do { written = write(child->pipes[0], data.data, data.len > 65536 ? 65536 : data.len); }
        while (written < 0 && errno == EINTR);
    int saved_errno = errno;
    if (!already_pending && sigpending(&pending) == 0 && sigismember(&pending, SIGPIPE) == 1) {
        int signal_number;
        sigwait(&blocked, &signal_number);
    }
    if (pthread_sigmask(SIG_SETMASK, &previous, NULL) != 0)
        return _TELIA_ERROR(StreamWriteFailed, "cannot restore signal mask");
    if (written >= 0) { *count = (uint64_t)written; return (_Error){0}; }
    if (saved_errno == EAGAIN || saved_errno == EWOULDBLOCK) { *would_block = true; return (_Error){0}; }
    return _TELIA_ERROR(StreamWriteFailed, "stdin pipe write failed");
#endif
}
