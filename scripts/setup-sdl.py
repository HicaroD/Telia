#!/usr/bin/env python3
"""Explicit, checksum-verified native SDL3 build. Never called by telia build."""
import argparse
import hashlib
import os
from pathlib import Path
import platform
import shutil
import subprocess
import tarfile
import urllib.request

VERSION = "3.2.28"
SHA256 = "1330671214d146f8aeb1ed399fc3e081873cdb38b5189d1f8bb6ab15bbc04211"
URL = "https://github.com/libsdl-org/SDL/releases/download/release-{0}/SDL3-{0}.tar.gz".format(VERSION)
ROOT = Path(__file__).resolve().parent.parent


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--archive", type=Path, help="use a local pinned source archive (no download)")
    parser.add_argument("--cache", type=Path, default=ROOT / ".native-cache")
    args = parser.parse_args()
    system = {"Darwin": "darwin", "Linux": "linux", "Windows": "windows"}.get(platform.system())
    arch = {"x86_64": "amd64", "AMD64": "amd64", "arm64": "arm64", "aarch64": "arm64"}.get(platform.machine())
    target = "{}-{}".format(system, arch)
    if target not in ("linux-amd64", "darwin-amd64", "darwin-arm64", "windows-amd64"):
        parser.error("unsupported native target: " + target)
    for tool in ("cmake", "ninja"):
        if not shutil.which(tool):
            parser.error("install {} and put it on PATH before setup".format(tool))
    cache = args.cache.resolve()
    cache.mkdir(parents=True, exist_ok=True)
    archive = args.archive or cache / ("SDL3-" + VERSION + ".tar.gz")
    if not archive.exists():
        if args.archive:
            parser.error("archive does not exist: " + str(archive))
        print("Downloading pinned SDL3 source: " + URL, flush=True)
        urllib.request.urlretrieve(URL, archive)
    if hashlib.sha256(archive.read_bytes()).hexdigest() != SHA256:
        parser.error("SDL3 source checksum mismatch; remove or replace " + str(archive))
    source = cache / ("SDL3-" + VERSION)
    if not source.exists():
        with tarfile.open(archive) as tar:
            # Validate paths and links before extracting a verified archive.
            for member in tar.getmembers():
                dest = (cache / member.name).resolve()
                if cache not in dest.parents or member.issym() or member.islnk():
                    parser.error("unsafe source archive member: " + member.name)
            tar.extractall(cache)
    prefix = ROOT / "base" / "vendor" / "sdl" / "artifacts" / target
    build = cache / ("build-" + VERSION + "-" + target)
    options = [
        "-DCMAKE_BUILD_TYPE=Release", "-DCMAKE_INSTALL_PREFIX=" + str(prefix),
        "-DCMAKE_INSTALL_LIBDIR=lib", "-DSDL_SHARED=ON", "-DSDL_STATIC=OFF",
        "-DSDL_TESTS=OFF", "-DSDL_EXAMPLES=OFF", "-DSDL_TEST_LIBRARY=OFF",
    ]
    if system == "darwin":
        options += ["-DCMAKE_INSTALL_NAME_DIR=@rpath", "-DCMAKE_OSX_ARCHITECTURES=" + ("arm64" if arch == "arm64" else "x86_64"), "-DCMAKE_OSX_DEPLOYMENT_TARGET=11.0"]
    if system == "linux":
        options += ["-DSDL_X11=ON", "-DSDL_WAYLAND=OFF"]
    if system == "windows":
        options += ["-DCMAKE_C_COMPILER=gcc", "-DSDL_LIBC=ON"]
    subprocess.run(["cmake", "-S", str(source), "-B", str(build), "-G", "Ninja"] + options, check=True)
    subprocess.run(["cmake", "--build", str(build), "--parallel", str(min(os.cpu_count() or 2, 8))], check=True)
    subprocess.run(["cmake", "--install", str(build)], check=True)
    shutil.copyfile(source / "LICENSE.txt", prefix / "LICENSE-SDL.txt")
    print("SDL3 ready for offline compilation: " + str(prefix))


if __name__ == "__main__":
    main()
