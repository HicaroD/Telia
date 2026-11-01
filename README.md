# The Telia Programming Language

Telia is a general-purpose programming language. Heavily inspired on the Go
Programming Language, Telia is simple in its essence.

```rust
package main

use "std::io"

fn main() i32 {
    io::println("Hello, Telia!")
    return 0
}
```

## Install

Available for Linux x64, macOS Intel/Apple Silicon, and Windows x64.
**A C compiler is required:**

- macOS: `xcode-select --install`
- Debian/Ubuntu: `sudo apt install build-essential`
- Windows: install [MSYS2](https://www.msys2.org/), then run
  `pacman -S mingw-w64-ucrt-x86_64-gcc` in its UCRT64 terminal and add
  `C:\msys64\ucrt64\bin` to your Windows `PATH`.

### macOS / Linux

```sh
curl -fsSL https://github.com/HicaroD/Telia/releases/latest/download/install.sh | sh
```

Add `~/.local/bin` to your shell profile's `PATH` for future terminals.

### Windows (PowerShell)

```powershell
irm https://github.com/HicaroD/Telia/releases/latest/download/install.ps1 | iex
```

Add `$env:LOCALAPPDATA\Telia\bin` to your user `PATH` for future terminals.

See the [installation guide](docs/getting-started/installation.md) for manual
installation, customization, and upgrades.

## Documentation

Start at the [documentation home](docs/README.md), or choose a path:

- **Learn:** [First program](docs/getting-started/first-program.md) and
  [language guide](docs/guide/README.md).
- **Look up:** [Language reference](docs/reference/README.md) and
  [standard library](docs/stdlib/README.md).
- **Build something:** [Cookbook](docs/cookbook/README.md).
- **Use an agent:** [Agent entry point](docs/agents.md).
- **Contribute:** [Development and documentation maintenance](docs/contributing/README.md).
