// OS-only entropy. Override at the native build boundary for deterministic
// tests, never via an application-level random/hasher callback.
#if defined(__linux__) || defined(__APPLE__)
#include <sys/random.h>
#elif defined(_WIN32)
#include <bcrypt.h>
#endif

static inline int _telia_os_entropy(unsigned char *bytes, size_t n) {
#if defined(__linux__)
    size_t done = 0;
    while (done < n) {
        ssize_t count = getrandom(bytes + done, n - done, 0);
        if (count < 0 && errno == EINTR) continue;
        if (count <= 0) return -1;
        done += (size_t)count;
    }
    return 0;
#elif defined(__APPLE__)
    // Our sole caller requests sixteen bytes (getentropy's limit is 256).
    int result;
    do { result = getentropy(bytes, n); } while (result < 0 && errno == EINTR);
    return result;
#elif defined(_WIN32)
    // Resolve the system provider without imposing a new application link flag.
    HMODULE module = LoadLibraryExW(L"bcrypt.dll", NULL, LOAD_LIBRARY_SEARCH_SYSTEM32);
    if (!module) return -1;
    typedef NTSTATUS (WINAPI *Generate)(BCRYPT_ALG_HANDLE, PUCHAR, ULONG, ULONG);
    Generate generate = (Generate)GetProcAddress(module, "BCryptGenRandom");
    int result = -1;
    if (generate && n <= ULONG_MAX &&
        generate(NULL, bytes, (ULONG)n, BCRYPT_USE_SYSTEM_PREFERRED_RNG) >= 0) result = 0;
    FreeLibrary(module);
    return result;
#else
    (void)bytes; (void)n;
    return -1; // Unsupported systems fail closed, never select a weak seed.
#endif
}

#ifndef _TELIA_ENTROPY
#define _TELIA_ENTROPY _telia_os_entropy
#endif

static inline _Error _telia_hash_seed(uint64_t *k0, uint64_t *k1) {
    unsigned char bytes[16] = {0};
    *k0 = 0; *k1 = 0;
    if (_TELIA_ENTROPY(bytes, sizeof(bytes)) != 0)
        return _TELIA_ERROR(RandomnessFailed, "OS entropy acquisition failed");
    for (unsigned i = 0; i < 8; ++i) {
        *k0 |= (uint64_t)bytes[i] << (8*i);
        *k1 |= (uint64_t)bytes[i+8] << (8*i);
    }
    return (_Error){0};
}
