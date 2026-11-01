#ifndef TELIA_TEST_EQUALITY_NATIVE_H
#define TELIA_TEST_EQUALITY_NATIVE_H
#include <math.h>
static double equality_nan(void) { return NAN; }
static double equality_infinity(void) { return INFINITY; }
static double equality_negative_zero(void) { return -0.0; }
static float equality_nan32(void) { return NAN; }
static float equality_negative_zero32(void) { return -0.0f; }
static const char *equality_malformed(void) { return "\xff\xfe"; }
#endif
