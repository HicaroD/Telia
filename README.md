# The Telia Programming Language

Telia is a general-purpose programming language. Heavily inspired on the Go
Programming Language, Telia is simple in its essence.

```rust
package main

use "std::io"

fn main() i32 {
    io::println("Hello, world!")
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

Download and inspect the installer before running it:

```sh
curl -fL https://github.com/HicaroD/Telia/releases/latest/download/install.sh -o install-telia.sh
sh install-telia.sh
export PATH="$HOME/.local/bin:$PATH"
telia version
```

Add `~/.local/bin` to your shell profile's `PATH` for future terminals.

### Windows (PowerShell)

```powershell
Invoke-WebRequest https://github.com/HicaroD/Telia/releases/latest/download/install.ps1 -OutFile install-telia.ps1
# Inspect the script before running it.
.\install-telia.ps1
$env:PATH = "$env:LOCALAPPDATA\Telia\bin;$env:PATH"
telia version
```

Add `$env:LOCALAPPDATA\Telia\bin` to your user `PATH` for future terminals.

See the [installation guide](docs/installation.md) for manual installation,
customization, upgrades, and release publishing.
