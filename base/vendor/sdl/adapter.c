#define SDL_MAIN_HANDLED
#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>

bool telia_sdl_init(void) {
    SDL_SetMainReady();
    return SDL_Init(SDL_INIT_VIDEO);
}
void telia_sdl_quit(void) { SDL_Quit(); }
const char *telia_sdl_get_error(void) { return SDL_GetError(); }
void *telia_sdl_create_window(const char *title, int32_t width, int32_t height) {
    return SDL_CreateWindow(title, width, height, 0);
}
void telia_sdl_destroy_window(void *window) { SDL_DestroyWindow(window); }
bool telia_sdl_raise_window(void *window) { return SDL_RaiseWindow(window); }
void *telia_sdl_create_renderer(void *window) { return SDL_CreateRenderer(window, NULL); }
void telia_sdl_destroy_renderer(void *renderer) { SDL_DestroyRenderer(renderer); }
uint32_t telia_sdl_poll_event(void) {
    SDL_Event event;
    return SDL_PollEvent(&event) ? event.type : 0;
}
bool telia_sdl_is_quit_event(uint32_t event) { return event == SDL_EVENT_QUIT; }
bool telia_sdl_set_draw_color(void *renderer, uint8_t r, uint8_t g, uint8_t b, uint8_t a) {
    return SDL_SetRenderDrawColor(renderer, r, g, b, a);
}
bool telia_sdl_clear(void *renderer) { return SDL_RenderClear(renderer); }
bool telia_sdl_fill_rect(void *renderer, float x, float y, float width, float height) {
    SDL_FRect rect = { x, y, width, height };
    return SDL_RenderFillRect(renderer, &rect);
}
bool telia_sdl_present(void *renderer) { return SDL_RenderPresent(renderer); }
void telia_sdl_delay(uint32_t milliseconds) { SDL_Delay(milliseconds); }
