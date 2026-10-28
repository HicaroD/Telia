#include "opaque.h"
struct FixtureWindow { uint32_t value; };
static FixtureWindow window;
FixtureWindow *fixture_window_create(uint32_t value) {
    window.value = value;
    return &window;
}
uint32_t fixture_window_value(FixtureWindow *value) { return value->value; }
void fixture_window_destroy(FixtureWindow *value) { value->value = 0; }
