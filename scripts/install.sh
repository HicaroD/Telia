#!/bin/sh
# Install a complete Telia release without sudo or shell-profile changes.
set -eu

fail() { echo "Error: $*" >&2; exit 1; }

if [ -t 1 ]; then
    cat <<'BANNER'
_________ _______  _       _________ _______
\__   __/(  ____ \( \      \__   __/(  ___  )
   ) (   | (    \/| (         ) (   | (   ) |
   | |   | (__    | |         | |   | (___) |
   | |   |  __)   | |         | |   |  ___  |
   | |   | (      | |         | |   | (   ) |
   | |   | (____/\| (____/\___) (___| )   ( |
   )_(   (_______/(_______/\_______/|/     \|

          THE TELIA PROGRAMMING LANGUAGE

BANNER
fi

for tool in curl tar mktemp; do
    command -v "$tool" >/dev/null 2>&1 || fail "Required tool not found: $tool"
done
case "$(uname -s)" in
    Darwin) os=darwin ;;
    Linux) os=linux ;;
    *) fail "Supported systems: macOS and Linux. Use install.ps1 on Windows." ;;
esac
case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    arm64|aarch64) arch=arm64 ;;
    *) fail "Unsupported CPU architecture: $(uname -m)" ;;
esac
target="$os-$arch"
case "$target" in
    darwin-amd64|darwin-arm64|linux-amd64) ;;
    *) fail "No release is available for $target" ;;
esac

repo=https://github.com/HicaroD/Telia
version=${TELIA_VERSION:-latest}
if [ "$version" = latest ]; then
    url=$(curl -fsSL --proto '=https' --proto-redir '=https' -o /dev/null -w '%{url_effective}' "$repo/releases/latest")
    version=${url##*/}
fi
case "$version" in
    v[0-9]*) ;;
    *) fail "Expected a version tag such as v0.1.0, got: $version" ;;
esac
case "$version" in *[!A-Za-z0-9.+-]*) fail "Invalid version tag" ;; esac

root=${TELIA_INSTALL_ROOT:-"$HOME/.local/share/telia"}
bin=${TELIA_BIN_DIR:-"$HOME/.local/bin"}
case "$root:$bin" in /*:/*) ;; *) fail "Installation directories must be absolute paths" ;; esac
destination="$root/$version"
[ ! -e "$destination" ] || fail "$destination already exists; choose another version or TELIA_INSTALL_ROOT"
if [ -e "$bin/telia" ] || [ -L "$bin/telia" ]; then
    [ -L "$bin/telia" ] || fail "Refusing to replace existing file: $bin/telia"
    case "$(readlink "$bin/telia")" in
        "$root"/v*/telia) ;;
        *) fail "Refusing to replace an unrelated symlink: $bin/telia" ;;
    esac
fi
mkdir -p "$root" "$bin"
work=$(mktemp -d "$root/.install.XXXXXX")
trap 'rm -rf "$work"' EXIT HUP INT TERM
asset="telia-$target.tar.gz"
base="$repo/releases/download/$version"
echo "Installing Telia $version for $target..."
curl -fsSL --proto '=https' --proto-redir '=https' "$base/$asset" -o "$work/$asset"
curl -fsSL --proto '=https' --proto-redir '=https' "$base/SHA256SUMS" -o "$work/SHA256SUMS"
expected=$(awk -v name="$asset" '$2 == name { print $1 }' "$work/SHA256SUMS")
[ "${#expected}" -eq 64 ] || fail "Missing or invalid checksum for $asset"
case "$expected" in *[!0-9a-fA-F]*) fail "Invalid checksum" ;; esac
if command -v sha256sum >/dev/null 2>&1; then
    actual=$(sha256sum "$work/$asset")
elif command -v shasum >/dev/null 2>&1; then
    actual=$(shasum -a 256 "$work/$asset")
else
    fail "Install sha256sum or shasum to verify the download"
fi
[ "$expected" = "${actual%% *}" ] || fail "Checksum mismatch; nothing was installed"
mkdir "$work/bundle"
tar -xzf "$work/$asset" -C "$work/bundle"
[ -x "$work/bundle/telia" ] && [ -d "$work/bundle/base/std" ] || fail "Incomplete release bundle"
mv "$work/bundle" "$destination"
ln -sfn "$destination/telia" "$bin/telia"
echo "Installed in $destination"
echo "Add this directory to PATH: $bin"
echo "Then run: telia version"
if ! command -v cc >/dev/null 2>&1 && ! command -v gcc >/dev/null 2>&1 && ! command -v clang >/dev/null 2>&1; then
    echo "A C compiler is required: macOS: xcode-select --install; Debian/Ubuntu: sudo apt install build-essential" >&2
fi
