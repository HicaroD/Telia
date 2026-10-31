#!/usr/bin/env python3
"""Verify an extracted release works outside the source checkout."""
import argparse
import hashlib
import os
from pathlib import Path
import subprocess
import tempfile


def check_bundle(bundle, launcher, root, env, version):
    output = subprocess.check_output([str(launcher), "version"], cwd=root, env=env, text=True)
    if not output.startswith("telia " + version + " ("):
        raise RuntimeError("incorrect release version: " + output)
    applications = root / "applications with spaces"
    applications.mkdir()
    suffix = ".exe" if os.name == "nt" else ""
    program = applications / ("hello" + suffix)
    subprocess.run([str(launcher), "build", str(bundle / "examples" / "hello.t"),
                    "-o", str(program)], cwd=root, env=env, check=True)
    output = subprocess.check_output([str(program)], cwd=root, env=env, text=True)
    if output.strip() != "Hello, world!":
        raise RuntimeError("unexpected hello-world output: " + output)
    # Exercise std native sources as well as vendor facade/raw transitive imports.
    source = root / "native.t"
    source.write_text('''package main
use "std::strings"
use "std::terminal"
use "vendor::sdl"
fn main() i32 {
  output := terminal::stdout()
  text := strings::clone("installed std") @fail
  defer strings::free(&text)
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    return 1
  }
  defer sdl::quit()
  window := sdl::create_window("installed SDL", 64, 48, 0) @fail
  defer sdl::destroy_window(window)
  size := sdl::get_window_size(window) @fail
  if size.width != 64 or size.height != 48 {
    return 2
  }
  return 0
}
''')
    native = applications / ("native" + suffix)
    subprocess.run([str(launcher), "build", str(source), "-o", str(native)],
                   cwd=root, env=env, check=True)
    relocated = root / "relocated applications"
    applications.rename(relocated)
    # Hide the installation so an absolute link/load path cannot make this pass.
    hidden = bundle.with_name(bundle.name + "-hidden")
    bundle.rename(hidden)
    try:
        runtime_env = dict(env, SDL_VIDEO_DRIVER="dummy", SDL_VIDEODRIVER="dummy",
                           SDL_RENDER_DRIVER="software")
        runtime_env["PATH"] = (str(Path(os.environ["SystemRoot"]) / "System32")
                               if os.name == "nt" else "/usr/bin:/bin")
        subprocess.run([str(relocated / native.name)], cwd=root, env=runtime_env,
                       check=True, timeout=60)
    finally:
        hidden.rename(bundle)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("archive", type=Path, nargs="?")
    parser.add_argument("--installed-bundle", type=Path)
    parser.add_argument("--launcher", type=Path)
    parser.add_argument("--version", required=True)
    args = parser.parse_args()
    with tempfile.TemporaryDirectory(prefix="telia-smoke-") as temp:
        root = Path(temp)
        env = {key: value for key, value in os.environ.items()
               if key not in ("T_STD_PATH", "T_VENDOR_PATH", "LD_LIBRARY_PATH",
                              "LD_PRELOAD", "DYLD_LIBRARY_PATH", "DYLD_FALLBACK_LIBRARY_PATH",
                              "DYLD_INSERT_LIBRARIES")}
        env["XDG_CONFIG_HOME"] = str(root / "config")
        if args.installed_bundle:
            if not args.launcher:
                parser.error("--installed-bundle requires --launcher")
            bundle, launcher = args.installed_bundle.resolve(), args.launcher.resolve()
        elif not args.archive:
            parser.error("provide an archive or --installed-bundle")
        elif os.name == "nt":
            bundle = root / "bundle with spaces"
            bundle.mkdir()
            subprocess.run(["tar", "-xzf", str(args.archive.resolve()), "-C", str(bundle)], check=True)
            launcher = bundle / "telia.exe"
        else:
            # Run the real installer, replacing only its network transport.
            tools = root / "tools"
            tools.mkdir()
            archive = args.archive.resolve()
            sums = root / "SHA256SUMS"
            sums.write_text(hashlib.sha256(archive.read_bytes()).hexdigest() + "  " + archive.name + "\n")
            curl = tools / "curl"
            curl.write_text('#!/bin/sh\nwhile [ "$#" -gt 0 ]; do\n'
                            'case "$1" in https://*) asset=${1##*/};; -o) shift; output=$1;; esac\n'
                            'shift\ndone\ncase "$asset" in SHA256SUMS) cp "$SMOKE_SUMS" "$output";;\n'
                            '*) cp "$SMOKE_ARCHIVE" "$output";; esac\n')
            curl.chmod(0o755)
            install = root / "installation with spaces"
            install_env = dict(env, PATH=str(tools) + os.pathsep + env["PATH"],
                               TELIA_VERSION="v0.0.0-smoke", TELIA_INSTALL_ROOT=str(install),
                               TELIA_BIN_DIR=str(root / "bin"), SMOKE_ARCHIVE=str(archive),
                               SMOKE_SUMS=str(sums))
            subprocess.run(["sh", str(Path(__file__).with_name("install.sh"))],
                           cwd=root, env=install_env, check=True)
            bundle = install / "v0.0.0-smoke"
            launcher = root / "bin" / "telia"
        check_bundle(bundle, launcher, root, env, args.version)
    print("Release smoke test passed")


if __name__ == "__main__":
    main()
