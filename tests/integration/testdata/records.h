#ifndef TELIA_FIXTURE_RECORDS_H
#define TELIA_FIXTURE_RECORDS_H
#if defined(__aarch64__) && defined(__clang__)
#define FIXTURE_CC __attribute__((aarch64_vector_pcs))
#elif defined(__x86_64__) && (defined(__GNUC__) || defined(__clang__))
#define FIXTURE_CC __attribute__((ms_abi))
#else
#define FIXTURE_CC
#endif
#include <stdint.h>
#include <stdbool.h>
#include <stdarg.h>
typedef enum { FIXTURE_KIND_READY = 37 } FixtureKind;
FixtureKind fixture_kind(bool known);
#define FIXTURE_RECORD_FLAG (UINT64_C(1) << 40)
typedef struct {
    int32_t value;
    uint64_t header_only;
    int32_t omitted;
} FixtureRecord;
int32_t fixture_record_check(FixtureRecord *record);
FixtureRecord fixture_record_roundtrip(FixtureRecord record);
typedef union { int32_t integer; double real; } FixtureUnion;
int32_t fixture_union_check(FixtureUnion value);
typedef struct { uint8_t data[4]; } FixtureBytes;
int32_t fixture_bytes_check(FixtureBytes value, uint8_t *first);
int32_t fixture_text_check(const char *text);
void fixture_text_replace(const char **output);
int32_t fixture_text_inspect(const char * const *texts);
typedef int32_t (*FixtureCompare)(int32_t left, int32_t right);
int32_t fixture_callback_invoke(FixtureCompare callback);
void fixture_callback_get(FixtureCompare *output);
typedef struct { FixtureCompare compare; } FixtureCallbacks;
FixtureCompare fixture_callback_return(void);
typedef int32_t (FIXTURE_CC *FixtureConventionCompare)(int32_t, int32_t);
int32_t FIXTURE_CC fixture_convention_invoke(FixtureConventionCompare callback);
int32_t fixture_variadic_check(int32_t marker, ...);
#pragma pack(push, 1)
typedef struct { uint8_t tag; uint64_t value; } FixturePacked;
#pragma pack(pop)
int32_t fixture_packed_check(FixturePacked *record, uint64_t *field, uint64_t size, uint64_t alignment);
FixturePacked fixture_packed_roundtrip(FixturePacked record);
uint64_t fixture_long_size(void);
long fixture_long_increment(long *value);
const char *fixture_plain_char(void);
int32_t fixture_char_check(const char *value);
typedef int32_t (*FixtureInspect)(va_list args);
int32_t fixture_va_drive(FixtureInspect callback);
int32_t fixture_va_check(va_list args);
typedef int32_t (*FixtureVariadic)(int32_t marker, ...);
int32_t fixture_va_owned_drive(FixtureVariadic callback);
int32_t fixture_va_order_drive(FixtureInspect callback);
int32_t fixture_va_pointer_drive(FixtureInspect callback);
#endif
