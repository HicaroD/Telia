#!/usr/bin/env python3
"""Package a native Telia release with std and prepared SDL3 dependencies."""
import argparse
import json
from pathlib import Path
import platform
import re
import shutil
import subprocess
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parent.parent


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
    package = ROOT / "base" / "vendor" / "sdl"
    manifest = json.loads((package / "native.json").read_text())
    entry = manifest["targets"].get(target)
    if entry is None:
        parser.error("unsupported release target: " + target)
    for path in entry["include_dirs"] + entry["libraries"] + entry["runtime_files"] + entry.get("sources", []):
        if not (package / path).exists():
            parser.error("missing SDL3 artifact {}; run scripts/setup-sdl.py first".format(path))
    with tempfile.TemporaryDirectory(prefix="telia-release-") as temp:
        stage = Path(temp)
        name = "telia.exe" if system == "windows" else "telia"
        flags = "-X main.Version={} -X main.Commit={}".format(args.version, args.commit)
        subprocess.run(["go", "build", "-trimpath", "-ldflags", flags, "-o", str(stage / name), "./cmd/compiler"], cwd=ROOT, check=True)
        shutil.copytree(ROOT / "base" / "std", stage / "base" / "std")
        dest = stage / "base" / "vendor" / "sdl"
        dest.mkdir(parents=True)
        for path in package.iterdir():
            if path.is_file():
                shutil.copy2(path, dest / path.name)
            elif path.name != "artifacts":
                shutil.copytree(path, dest / path.name)
        shutil.copytree(package / "artifacts" / target, dest / "artifacts" / target)
        shutil.copy2(ROOT / "LICENSE", stage / "LICENSE")
        shutil.copy2(ROOT / "docs" / "native-libraries.md", stage / "native-libraries.md")
        examples = stage / "examples"
        examples.mkdir()
        shutil.copy2(ROOT / "examples" / "sdl.t", examples / "sdl.t")
        shutil.copy2(ROOT / "examples" / "hello.t", examples / "hello.t")
        args.output.parent.mkdir(parents=True, exist_ok=True)
        with tarfile.open(args.output, "w:gz", dereference=True) as archive:
            for path in sorted(stage.iterdir()):
                archive.add(path, arcname=path.name)
    print("Ready-to-use {} release: {}".format(target, args.output))


if __name__ == "__main__":
    main()
