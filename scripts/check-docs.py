#!/usr/bin/env python3
"""Check local Markdown links/heading anchors without network dependencies."""

import re
import sys
from pathlib import Path
from urllib.parse import unquote, urlsplit


ROOT = Path(__file__).resolve().parents[1]
LINK = re.compile(r"!?\[[^\]\n]*\]\(([^\s)]+)(?:\s+\"[^\"]*\")?\)")


def prose(text):
    """Discard fenced examples so their sample links/headings are not checked."""
    lines = []
    fence = None
    for line in text.splitlines():
        marker = re.match(r"^\s*(`{3,}|~{3,})", line)
        if marker:
            run = marker.group(1)
            if fence is None:
                fence = run
            elif run[0] == fence[0] and len(run) >= len(fence):
                fence = None
            lines.append("")
        else:
            lines.append(line if fence is None else "")
    return "\n".join(lines)


def anchors(text):
    result = set()
    counts = {}
    for line in prose(text).splitlines():
        heading = re.match(r"^#{1,6}\s+(.+?)\s*#*\s*$", line)
        if heading:
            title = re.sub(r"\[([^]]+)\]\([^)]+\)", r"\1", heading.group(1))
            slug = re.sub(r"[^\w\- ]", "", title.lower()).replace(" ", "-")
            count = counts.get(slug, 0)
            counts[slug] = count + 1
            result.add(slug if count == 0 else f"{slug}-{count}")
        result.update(re.findall(r'<a\s+(?:id|name)="([^"]+)"', line))
    return result


def check(paths, root=ROOT):
    failures = []
    cache = {}
    for path in sorted(paths):
        text = prose(path.read_text(encoding="utf-8"))
        for number, line in enumerate(text.splitlines(), 1):
            for match in LINK.finditer(line):
                link = match.group(1).strip("<>")
                url = urlsplit(link)
                if url.scheme or url.netloc:
                    continue
                target = (path.parent / unquote(url.path)).resolve() if url.path else path.resolve()
                label = f"{path.relative_to(root)}:{number}: {link}"
                if not target.exists():
                    failures.append(f"{label}: missing target")
                elif url.fragment and target.suffix.lower() == ".md":
                    if target not in cache:
                        cache[target] = anchors(target.read_text(encoding="utf-8"))
                    if unquote(url.fragment) not in cache[target]:
                        failures.append(f"{label}: missing heading anchor")
    return failures


def main():
    paths = {ROOT / "README.md", ROOT / "AGENTS.md"}
    paths.update((ROOT / "docs").rglob("*.md"))
    paths.update((ROOT / "examples").rglob("*.md"))
    failures = check(paths)
    failures.extend(check_packages(ROOT))
    for failure in failures:
        print(failure, file=sys.stderr)
    if failures:
        return 1
    print(f"Documentation links OK ({len(paths)} Markdown files)")
    return 0


def check_packages(root):
    """Each standard-library source package must be discoverable in its index."""
    source = root / "base" / "std"
    packages = {
        "std::" + "::".join(path.parent.relative_to(source).parts)
        for path in source.rglob("*.t")
    }
    index = (root / "docs" / "stdlib" / "README.md").read_text(encoding="utf-8")
    listed = set(re.findall(r"\| `(std::[\w:]+)` \|", index))
    return [f"docs/stdlib/README.md: missing package {name}" for name in sorted(packages - listed)]


if __name__ == "__main__":
    sys.exit(main())
