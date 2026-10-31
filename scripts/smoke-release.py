#!/usr/bin/env python3
"""Verify an extracted release works outside the source checkout."""
import argparse
import os
from pathlib import Path
import subprocess
import tempfile


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("archive", type=Path)
    parser.add_argument("--version", required=True)
    args = parser.parse_args()
    with tempfile.TemporaryDirectory(prefix="telia-smoke-") as temp:
        root = Path(temp)
        bundle = root / "bundle"
        bundle.mkdir()
        subprocess.run(["tar", "-xzf", str(args.archive.resolve()), "-C", str(bundle)], check=True)
        env = {key: value for key, value in os.environ.items()
               if key not in ("T_STD_PATH", "T_VENDOR_PATH")}
        env["XDG_CONFIG_HOME"] = str(root / "config")
        binary = bundle / ("telia.exe" if os.name == "nt" else "telia")
        version = subprocess.check_output([str(binary), "version"], cwd=root, env=env, text=True)
        if not version.startswith("telia " + args.version + " ("):
            raise RuntimeError("incorrect release version: " + version)
        program = root / ("hello.exe" if os.name == "nt" else "hello")
        subprocess.run([str(binary), "build", str(bundle / "examples" / "hello.t"),
                        "-o", str(program)], cwd=root, env=env, check=True)
        output = subprocess.check_output([str(program)], cwd=root, env=env, text=True)
        if output.strip() != "Hello, world!":
            raise RuntimeError("unexpected hello-world output: " + output)
    print("Release smoke test passed")


if __name__ == "__main__":
    main()
