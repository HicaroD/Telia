# The Telia Programming Language

Telia is a general-purpose programming language. Heavily inspired on the Go
Programming Language, Telia is simple in its essence.

```rust
package main

use "std::io"
use "std::strings"

fn main() i32 {
    greeting := strings::clone("Hello, Telia!") @fail
    defer strings::free(&greeting)

    io::println(strings::view(&greeting))
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

See the [installation guide](docs/installation.md) for manual installation,
customization, upgrades, and release publishing.

## Language guides

- [Unit testing](docs/testing.md)
- [Collections and composed field access](docs/collections.md)
- [Explicit integer and pointer conversions](docs/conversions.md)
- [Errors and cleanup-safe propagation](docs/errors.md)
- [Memory ownership and transfer](docs/strings-and-memory.md#ownership-conventions-and-transfer)
