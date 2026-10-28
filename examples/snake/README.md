# Snake in Telia + SDL3

A 28 × 24 Snake game written in Telia. No image, font, or SDL extension libraries
are needed: the interface uses SDL3's built-in debug text renderer.

From the repository root, after [SDL setup](../../docs/native-libraries.md):

```sh
go build -o telia ./cmd/compiler
./telia build examples/snake -o bin/snake
./bin/snake
```

Keep the deployed SDL shared library beside the executable.

| Control | Action |
| --- | --- |
| Arrow keys / WASD | Turn |
| Space | Pause / resume |
| R | Restart |
| Escape / close window | Quit |

Eat the orange food to grow. Hitting a wall or your body ends the game. Reverse
turns are ignored, and only one turn is accepted per movement tick. Speed rises
at scores 10 and 20. Filling the board wins. The session best survives restarts
but is not saved to disk.

`game/game.t` contains deterministic game rules with no SDL dependency. The UI,
input, and fixed-step clock live in `main.t`.

For an automated, invisible smoke run:

```sh
SDL_VIDEO_DRIVER=dummy SDL_RENDER_DRIVER=software TELIA_SNAKE_FRAMES=12 ./bin/snake
```

Leave these variables unset when playing normally.
