#!/usr/bin/env python3
"""Report formatting differences in tracked Go files without changing them."""

import subprocess
import sys


def main():
    paths = subprocess.check_output(
        ["git", "ls-files", "-z", "--", "*.go"]
    ).decode().split("\0")
    paths = [path for path in paths if path]
    if not paths:
        return 0

    different = False
    for command in (
        ["gofmt", "-d"],
        ["golines", "--base-formatter=gofmt", "--dry-run"],
    ):
        result = subprocess.run(command + paths, stdout=subprocess.PIPE, check=True)
        if result.stdout:
            sys.stdout.buffer.write(result.stdout)
            different = True
    if different:
        print("Formatting differences found; run make fmt.", file=sys.stderr)
    return int(different)


if __name__ == "__main__":
    sys.exit(main())
