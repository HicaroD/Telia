#include <SDL3/SDL.h>

/* Keep SDL's allocator paired with its deallocator without pointer casts in
 * application code or incompatible redeclarations of SDL_free. */
void telia_sdl_free_text(char *text) {
    SDL_free(text);
}

bool telia_sdl_render_debug_i32(SDL_Renderer *renderer, float x, float y, Sint32 value) {
    return SDL_RenderDebugTextFormat(renderer, x, y, "%" SDL_PRIs32, value);
}

typedef Sint32 (*TeliaThreadEntry)(void *);
typedef struct {
    TeliaThreadEntry entry;
    void *data;
} TeliaThreadStart;

static int SDLCALL telia_thread_entry(void *opaque) {
    TeliaThreadStart *start = opaque;
    TeliaThreadEntry entry = start->entry;
    void *data = start->data;
    SDL_free(start);
    return entry(data);
}

SDL_Thread *telia_sdl_create_thread(TeliaThreadEntry entry, const char *name, void *data) {
    if (!entry) {
        SDL_SetError("Missing Telia thread entry");
        return NULL;
    }
    TeliaThreadStart *start = SDL_malloc(sizeof(*start));
    if (!start) {
        SDL_OutOfMemory();
        return NULL;
    }
    start->entry = entry;
    start->data = data;
    SDL_Thread *thread = SDL_CreateThread(telia_thread_entry, name, start);
    if (!thread) {
        SDL_free(start);
    }
    return thread;
}

SDL_IOStream *telia_sdl_io_from_bytes(Uint8 *bytes, Uint64 size) {
    return SDL_IOFromMem(bytes, (size_t)size);
}

Uint64 telia_sdl_read_bytes(SDL_IOStream *stream, Uint8 *bytes, Uint64 size) {
    return SDL_ReadIO(stream, bytes, (size_t)size);
}
