# Installation and releases

For quick-install commands and C compiler setup, see the [README](../README.md#install).
Releases bundle the compiler, standard library, and SDL3 package. Go is not
required. Both installers verify SHA-256 checksums, require no administrator
privileges, and leave shell profiles unchanged.

## Installation locations

| Platform | Bundle | Command on `PATH` |
| --- | --- | --- |
| macOS / Linux | `~/.local/share/telia/<version>` | `~/.local/bin/telia` (symlink) |
| Windows | `%LOCALAPPDATA%\Telia\<version>` | `%LOCALAPPDATA%\Telia\bin\telia.cmd` (launcher) |

## Customization

Pin a version or choose absolute installation paths:

```sh
TELIA_VERSION=v0.1.0 TELIA_INSTALL_ROOT="$HOME/tools/telia" TELIA_BIN_DIR="$HOME/tools/bin" sh install-telia.sh
```

```powershell
.\install-telia.ps1 -Version v0.1.0 -InstallRoot C:\Tools\Telia
```

The default is the latest stable release. Installing an already present version
is refused. Custom command directories must be added to `PATH` separately.

## Manual installation and troubleshooting

Download the matching archive and `SHA256SUMS` from
[GitHub Releases](https://github.com/HicaroD/Telia/releases), verify the checksum,
and extract it. Add the extracted directory to `PATH`. **Keep the executable and
`base/` together**; copying only the executable breaks library discovery.

If PowerShell blocks scripts, follow your organization's execution policy or use
manual installation. Existing `T_STD_PATH` / `T_VENDOR_PATH` settings override
bundled libraries; unset them when using the release bundle.

## Upgrades and removal

Run the installer again after a new release to upgrade. Previous version folders
are kept. To uninstall, delete the installation folder and installer-created
symlink/launcher, then remove any `PATH` entry you added.

## Publishing a release

Push a version tag pointing to the commit you want to ship:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The **Native vendor packages** workflow builds and tests all four platforms,
extracts each bundle outside the checkout, compiles and runs hello-world, then
publishes the archives, installers, and `SHA256SUMS` together. Tags containing
`-` (such as `v0.1.0-rc.1`) are published as prereleases.

Until the first stable release is published, release-based installers will not
work. Bundles are also available from successful workflow runs under **Actions**.
