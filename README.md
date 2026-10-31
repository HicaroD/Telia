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

> [!WARNING]
> Telia compiler is heavily imature and not ready to use.

## Install

Prebuilt releases include the compiler, standard library, and SDL3 package for
Linux x64, macOS Intel/Apple Silicon, and Windows x64. Go is not required.
**You still need a C compiler** to compile Telia programs:

- macOS: `xcode-select --install`
- Debian/Ubuntu: `sudo apt install build-essential`
- Windows: install [MSYS2](https://www.msys2.org/), then run
  `pacman -S mingw-w64-ucrt-x86_64-gcc` in its UCRT64 terminal and add
  `C:\msys64\ucrt64\bin` to your Windows `PATH`.

### macOS / Linux

Download and inspect the installer, then run it:

```sh
curl -fL https://github.com/HicaroD/Telia/releases/latest/download/install.sh -o install-telia.sh
sh install-telia.sh
export PATH="$HOME/.local/bin:$PATH"
telia version
```

The installer verifies the archive's SHA-256 checksum and installs the complete
bundle under `~/.local/share/telia/<version>`, with a symlink at
`~/.local/bin/telia`. Add `~/.local/bin` to your shell profile's `PATH` to keep it
available in new terminals.

To pin a version or customize directories:

```sh
TELIA_VERSION=v0.1.0 TELIA_INSTALL_ROOT="$HOME/tools/telia" TELIA_BIN_DIR="$HOME/tools/bin" sh install-telia.sh
```

### Windows (PowerShell)

```powershell
Invoke-WebRequest https://github.com/HicaroD/Telia/releases/latest/download/install.ps1 -OutFile install-telia.ps1
# Inspect the script before running it.
.\install-telia.ps1
$env:PATH = "$env:LOCALAPPDATA\Telia\bin;$env:PATH"
telia version
```

The installer verifies the checksum and installs under
`$env:LOCALAPPDATA\Telia\<version>`. Add `$env:LOCALAPPDATA\Telia\bin` to your user
`PATH`, then open a new terminal. It contains a launcher for the selected version.
Use `-Version v0.1.0` or `-InstallRoot C:\Tools\Telia` to customize the installation.
If your execution policy blocks scripts, follow your organization's policy;
manual installation below does not require running the installer.

### Manual installation, upgrades, and removal

You can also download the matching archive and `SHA256SUMS` from
[GitHub Releases](https://github.com/HicaroD/Telia/releases), verify its checksum,
and extract it. Add the extracted directory to `PATH`. **Keep the executable and
`base/` together**; copying only the executable breaks library discovery.

Run the installer again after a new release to upgrade. Previous version folders
are kept; installing an already present version is refused. To remove Telia,
delete its installation folder and the installer-created symlink/launcher, and
remove any `PATH` entry you added. Neither installer requires administrator
privileges or edits shell profiles. Existing `T_STD_PATH` / `T_VENDOR_PATH`
settings override bundled libraries; unset these when using the release bundle.

## Publishing a release

Push a version tag pointing to the commit you want to ship:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The **Native vendor packages** workflow builds and tests all four platforms,
extracts each bundle outside the checkout, compiles and runs hello-world, then
publishes the archives, installers, and `SHA256SUMS` together. Tags containing
`-` (such as `v0.1.0-rc.1`) are published as prereleases. The default installer
selects the latest stable release. Until the first release is published, bundles
are available from successful workflow runs under **Actions**, but release-based
installers will not work yet.
