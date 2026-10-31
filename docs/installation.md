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

Push your changes to `master` and wait for both **CI** and **Native vendor
packages** to succeed for the exact commit you want to ship. Then tag that
tested commit and push the tag:

```sh
git tag v0.2.1 TESTED_COMMIT_SHA
git push origin v0.2.1
```

The tag starts the lightweight **Release** workflow, not another full CI/native
test run. It verifies successful `master` push runs of both workflows for that
exact commit and downloads the tested SDL artifacts. It then rebuilds the
compiler with the release version embedded, packages all four platforms, and
smoke-tests the extracted bundles outside the checkout. Windows also checks its
installer. Only when all platforms pass are the archives, installers, and
`SHA256SUMS` published together. Tags containing `-` (such as `v0.2.1-rc.1`) are
published as prereleases. Published releases cannot be overwritten.

Release stops if either workflow is unfinished or unsuccessful, or if its
matching native artifacts are missing, expired, or ambiguous. It never silently
rebuilds SDL or borrows artifacts from another commit, a PR, or a fork. Native
artifacts are retained for 30 days; rerun the original **Native vendor packages**
`master` push run for that commit if they expire, then rerun the failed release.
Fix a failing check and test the new commit before creating a new version tag.

Development bundles are also available for seven days from successful native
workflow runs under **Actions**.
