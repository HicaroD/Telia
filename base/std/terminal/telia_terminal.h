#ifndef TELIA_TERMINAL_H
#define TELIA_TERMINAL_H
#include <stdint.h>
#include <stdbool.h>

intptr_t telia_terminal_standard(int32_t stream);
int32_t telia_terminal_close(intptr_t handle);
int32_t telia_terminal_open(bool input, intptr_t *handle);
int32_t telia_terminal_write(intptr_t handle, void *bytes, uint64_t length, uint64_t *count);
int32_t telia_terminal_write_text(intptr_t handle, char *bytes, uint64_t length, uint64_t *count);
int32_t telia_terminal_read(intptr_t handle, void *bytes, uint64_t length,
                           void *pending_utf8, uint64_t *pending_count, uint32_t *pending_utf16,
                           uint64_t *count, bool *eof);
int32_t telia_terminal_is_terminal(intptr_t handle, bool *result);
int32_t telia_terminal_size(intptr_t handle, uint64_t *columns, uint64_t *rows);
int32_t telia_terminal_get_mode(intptr_t handle, bool *echo, bool *canonical, bool *signals);
int32_t telia_terminal_apply_mode(intptr_t handle, bool raw, bool echo, bool canonical, bool signals, void **saved);
int32_t telia_terminal_restore(void **saved);
int32_t telia_terminal_prepare_output(intptr_t handle, void **saved);
int32_t telia_terminal_profile(intptr_t handle, bool *dumb, bool *ansi, bool *indexed, bool *rgb);
int32_t telia_terminal_wait(intptr_t handle, bool input, bool resize, int64_t timeout_ms, bool buffered_input,
                           bool *initialized, uint64_t *columns, uint64_t *rows,
                           bool *ready, bool *resized, bool *timed_out);
#endif
