#include "records.h"
#include <string.h>
#include <stdarg.h>
#include <stddef.h>
int32_t fixture_packed_check(FixturePacked *record, uint64_t *field, uint64_t size, uint64_t alignment) {
    return size == sizeof(FixturePacked) && alignment == _Alignof(FixturePacked) &&
        (char *)field - (char *)record == offsetof(FixturePacked, value) &&
        record->tag == 7 && record->value == 40;
}
FixturePacked fixture_packed_roundtrip(FixturePacked record) { record.value += 2; return record; }
uint64_t fixture_long_size(void) { return sizeof(long); }
long fixture_long_increment(long *value) { return ++*value; }
const char *fixture_plain_char(void) { return "hello"; }
int32_t fixture_char_check(const char *value) { return strcmp(value, "hello") == 0; }
int32_t fixture_va_check(va_list args) {
    int number = va_arg(args, int);
    double fraction = va_arg(args, double);
    return number == 42 && fraction == 1.5;
}
static int32_t fixture_va_call(FixtureInspect callback, int marker, ...) {
    va_list args;
    va_start(args, marker);
    int32_t result = callback(args);
    va_end(args);
    return result;
}
int32_t fixture_va_drive(FixtureInspect callback) { return fixture_va_call(callback, 0, 42, 1.5); }
int32_t fixture_va_owned_drive(FixtureVariadic callback) { return callback(0, 42, 1.5); }
int32_t fixture_va_order_drive(FixtureInspect callback) { return fixture_va_call(callback, 0, 42, 7, 99); }
int32_t fixture_va_pointer_drive(FixtureInspect callback) {
    int32_t value = 42, sentinel = 99;
    return fixture_va_call(callback, 0, &value, &sentinel);
}
int32_t fixture_variadic_check(int32_t marker, ...) {
    va_list args;
    va_start(args, marker);
    int truth = va_arg(args, int);
    int small = va_arg(args, int);
    int wide = va_arg(args, int);
    double fraction = va_arg(args, double);
    va_end(args);
    return marker == 42 && truth == 1 && small == -7 && wide == 60000 && fraction == 1.5;
}
int32_t fixture_text_check(const char *text) { return strcmp(text, "hello") == 0; }
void fixture_text_replace(const char **output) { *output = "after"; }
int32_t fixture_text_inspect(const char * const *texts) { return strcmp(*texts, "after") == 0; }
static int32_t fixture_add(int32_t left, int32_t right) { return left + right; }
int32_t fixture_callback_invoke(FixtureCompare callback) { return callback(44, 2); }
void fixture_callback_get(FixtureCompare *output) { *output = fixture_add; }
FixtureCompare fixture_callback_return(void) { return fixture_add; }
int32_t FIXTURE_CC fixture_convention_invoke(FixtureConventionCompare callback) { return callback(44, 2); }
FixtureKind fixture_kind(bool known) { return known ? FIXTURE_KIND_READY : (FixtureKind)99; }
int32_t fixture_record_check(FixtureRecord *record) {
    return record->value == 42 && record->header_only == 0 && record->omitted == 0;
}
FixtureRecord fixture_record_roundtrip(FixtureRecord record) {
    record.value += 2;
    record.omitted = 7;
    return record;
}
int32_t fixture_union_check(FixtureUnion value) { return value.integer == 42; }
int32_t fixture_bytes_check(FixtureBytes value, uint8_t *first) {
    return value.data[1] == 42 && first[0] == 10 && first[3] == 40;
}
