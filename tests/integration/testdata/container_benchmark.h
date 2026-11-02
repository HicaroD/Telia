#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <stddef.h>
#include <string.h>
#include <time.h>

// Account requested table bytes, excluding malloc's own metadata and stack keys.
typedef union { max_align_t alignment; size_t bytes; } bench_header;
static size_t bench_live, bench_peak;
static void *bench_malloc(size_t bytes) {
    bench_header *p = malloc(sizeof(*p) + bytes);
    if (!p) return NULL;
    p->bytes = bytes;
    bench_live += bytes;
    if (bench_live > bench_peak) bench_peak = bench_live;
    return p + 1;
}
static void bench_free(void *ptr) {
    if (!ptr) return;
    bench_header *p = (bench_header *)ptr - 1;
    bench_live -= p->bytes;
    free(p);
}
#define _TELIA_MALLOC bench_malloc
#define _TELIA_FREE bench_free
// Fixed benchmark seed only: excludes entropy timing and removes seed variance.
static int bench_entropy(void *ptr, size_t size) { memset(ptr, 0x71, size); return 0; }
#define _TELIA_ENTROPY bench_entropy
static const char *_bench_key(uint64_t index) {
    static char keys[8192][32];
    snprintf(keys[index], sizeof(keys[index]), "key-%08llu", (unsigned long long)index);
    return keys[index];
}
static clock_t bench_start;
static double bench_seconds[6];
static uint64_t bench_counts[6];
static void _bench_start(void) { bench_start = clock(); }
static void _bench_stop(uint64_t op, uint64_t count) {
    bench_seconds[op] += (double)(clock() - bench_start) / CLOCKS_PER_SEC;
    bench_counts[op] += count;
}
static void _bench_report(uint64_t size, uint64_t checksum) {
    static const char *names[] = {"insert", "lookup", "iteration", "remove", "churn-pair", "growth"};
    if (bench_live != 0 || checksum == 0) abort();
    for (unsigned op = 0; op < 6; ++op) {
        printf("%llu,%s,%.3f,%zu,%llu\n", (unsigned long long)size, names[op],
            bench_seconds[op] * 1e9 / bench_counts[op], bench_peak,
            (unsigned long long)checksum);
        bench_seconds[op] = 0;
        bench_counts[op] = 0;
    }
    bench_peak = 0;
}
