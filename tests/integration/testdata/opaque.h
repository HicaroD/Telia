#ifndef TELIA_OPAQUE_FIXTURE_H
#define TELIA_OPAQUE_FIXTURE_H
#include <stdint.h>
typedef struct FixtureWindow FixtureWindow;
FixtureWindow *fixture_window_create(uint32_t value);
uint32_t fixture_window_value(FixtureWindow *window);
void fixture_window_destroy(FixtureWindow *window);
#endif
