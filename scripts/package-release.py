#!/usr/bin/env python3
"""Package a native Telia release with std and prepared SDL3 dependencies."""
import argparse
import hashlib
import json
from pathlib import Path
import platform
import re
import shutil
import subprocess
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parent.parent


def copy_libraries(source, destination, target):
    """Ship every package, but only the selected target's native artifacts."""
    def ignore(directory, names):
        if Path(directory).name == "artifacts":
            return [name for name in names if name != target]
        return []

    shutil.copytree(source, destination, ignore=ignore)
    for manifest_path in destination.rglob("native.json"):
        manifest = json.loads(manifest_path.read_text())
        entry = manifest["targets"].get(target)
        if not manifest.get("version") or entry is None:
            raise ValueError("unsupported or unversioned native package: " + str(manifest_path))
        for key in ("include_dirs", "sources", "libraries", "runtime_files"):
            for name in entry.get(key, []):
                path = Path(name)
                if path.is_absolute() or ".." in path.parts or "\\" in name:
                    raise ValueError("native paths must be package-relative: " + name)
                artifact = manifest_path.parent / path
                valid = artifact.is_dir() if key == "include_dirs" else artifact.is_file()
                if not valid:
                    raise ValueError("missing native artifact for {}: {}".format(target, artifact))


def write_inventory(stage):
    entries = []
    for path in sorted(stage.rglob("*")):
        if path.is_file():
            name = path.relative_to(stage).as_posix()
            if any(char in name for char in "\r\n\\"):
                raise ValueError("unsupported bundle filename: " + name)
            entries.append(hashlib.sha256(path.read_bytes()).hexdigest() + "  " + name + "\n")
    (stage / "BUNDLE_SHA256SUMS").write_text("".join(entries), encoding="utf-8")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, required=True, help="output .tar.gz archive")
    parser.add_argument("--version", default="dev", help="release tag embedded in the compiler")
    parser.add_argument("--commit", default="unknown", help="source commit embedded in the compiler")
    args = parser.parse_args()
    for value in (args.version, args.commit):
        if not re.fullmatch(r"[A-Za-z0-9.+-]+", value):
            parser.error("version and commit must contain only letters, digits, dots, pluses, or hyphens")
    system = {"Darwin": "darwin", "Linux": "linux", "Windows": "windows"}.get(platform.system())
    arch = {"x86_64": "amd64", "AMD64": "amd64", "arm64": "arm64", "aarch64": "arm64"}.get(platform.machine())
    target = "{}-{}".format(system, arch)
    if target not in ("linux-amd64", "darwin-amd64", "darwin-arm64", "windows-amd64"):
        parser.error("unsupported release target: " + target)
    with tempfile.TemporaryDirectory(prefix="telia-release-") as temp:
        stage = Path(temp)
        name = "telia.exe" if system == "windows" else "telia"
        flags = "-X main.Version={} -X main.Commit={}".format(args.version, args.commit)
        subprocess.run(["go", "build", "-trimpath", "-ldflags", flags, "-o", str(stage / name), "./cmd/compiler"], cwd=ROOT, check=True)
        for library in ("std", "vendor"):
            copy_libraries(ROOT / "base" / library, stage / "base" / library, target)
        shutil.copy2(ROOT / "LICENSE", stage / "LICENSE")
        shutil.copy2(ROOT / "docs" / "native-libraries.md", stage / "native-libraries.md")
        examples = stage / "examples"
        examples.mkdir()
        shutil.copy2(ROOT / "examples" / "sdl.t", examples / "sdl.t")
        shutil.copy2(ROOT / "examples" / "hello.t", examples / "hello.t")
        write_inventory(stage)
        args.output.parent.mkdir(parents=True, exist_ok=True)
        with tarfile.open(args.output, "w:gz", dereference=True) as archive:
            for path in sorted(stage.iterdir()):
                archive.add(path, arcname=path.name)
    print("Ready-to-use {} release: {}".format(target, args.output))


if __name__ == "__main__":
    main()
