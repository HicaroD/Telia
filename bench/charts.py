#!/usr/bin/env python3
"""Render compiler phase charts as SVG using only the Python standard library."""

import argparse
import html
import json
import statistics
from pathlib import Path


def render(path):
    data = json.loads(path.read_text())
    groups = {}
    for sample in data["samples"]:
        groups.setdefault(sample["functions"], []).append(sample)
    width = 960
    height = 135 + 110 * len(groups)
    svg = [f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}" viewBox="0 0 {width} {height}">',
           '<rect width="100%" height="100%" fill="white"/>',
           '<g font-family="sans-serif" font-size="14" fill="#17202a">',
           '<text x="24" y="28" font-size="20">Telia compiler: where build time goes</text>']
    phases = [("parse_ms", "Parse", "#2874a6"), ("sema_ms", "Semantic analysis", "#148f77"),
              ("backend_ms", "C backend + toolchain", "#b9770e")]
    maximum = max(max(s["parse_ms"] + s["sema_ms"] + s["backend_ms"] for s in rows) for rows in groups.values())
    scale = 650 / max(maximum, 0.001)
    for index, (n, rows) in enumerate(sorted(groups.items())):
        y = 65 + index * 110
        svg.append(f'<text x="24" y="{y}">{n} functions</text>')
        x = 190
        for key, label, color in phases:
            median = statistics.median(s[key] for s in rows)
            length = median * scale
            svg.append(f'<rect x="{x:.2f}" y="{y-16}" width="{length:.2f}" height="24" fill="{color}"><title>{html.escape(label)}: {median:.3f} ms median</title></rect>')
            x += length
        total = [s["parse_ms"] + s["sema_ms"] + s["backend_ms"] for s in rows]
        tool = statistics.median(s["toolchain_ms"] for s in rows)
        allocations = statistics.median(s["allocated_bytes"] for s in rows) / 1024 / 1024
        svg.append(f'<text x="190" y="{y+32}">Total median {statistics.median(total):.2f} ms; range {min(total):.2f}–{max(total):.2f} ms</text>')
        svg.append(f'<text x="190" y="{y+54}">Toolchain median {tool:.2f} ms (subset of backend); Go allocations {allocations:.2f} MiB/build</text>')
    y = height - 45
    for index, (_, label, color) in enumerate(phases):
        x = 24 + index * 270
        svg.append(f'<rect x="{x}" y="{y-12}" width="14" height="14" fill="{color}"/><text x="{x+22}" y="{y}">{label}</text>')
    svg.append(f'<text x="24" y="{height-16}" font-size="12">{data["runs"]} warm-cache samples/size; bars stack phase medians, not median total. No CLI startup.</text>')
    svg.append('</g></svg>')
    path.with_name("phases.svg").write_text("\n".join(svg))


def compare(path, baseline):
    before = json.loads(baseline.read_text())
    after = json.loads(path.read_text())
    sizes = sorted({s["functions"] for s in before["samples"]} &
                   {s["functions"] for s in after["samples"]})
    height = 100 + 100 * len(sizes)
    svg = [f'<svg xmlns="http://www.w3.org/2000/svg" width="1000" height="{height}" viewBox="0 0 1000 {height}">',
           '<rect width="100%" height="100%" fill="white"/>',
           '<g font-family="sans-serif" font-size="13" fill="#17202a">',
           '<text x="24" y="28" font-size="20">Buffered lookahead: before / after</text>',
           '<text x="24" y="50">Gray: baseline. Green: updated. Medians; lower is better.</text>']
    for column, (key, title, divisor) in enumerate([
        ("parse_ms", "Parsing (ms)", 1),
        ("allocated_bytes", "Go allocation volume/build (MiB)", 1024 * 1024),
    ]):
        x = 130 + 470 * column
        svg.append(f'<text x="{x}" y="78">{title}</text>')
        values = {}
        for n in sizes:
            values[n] = [statistics.median(s[key] for s in data["samples"] if s["functions"] == n) / divisor
                         for data in (before, after)]
        maximum = max((max(v) for v in values.values()), default=1)
        scale = 270 / max(maximum, 0.001)
        for row, n in enumerate(sizes):
            y = 100 + row * 100
            if column == 0:
                svg.append(f'<text x="24" y="{y+15}">{n} functions</text>')
            for index, (value, color) in enumerate(zip(values[n], ("#85929e", "#148f77"))):
                bar_y = y + index * 27
                width = value * scale
                svg.append(f'<rect x="{x}" y="{bar_y}" width="{width:.2f}" height="20" fill="{color}"/>')
                svg.append(f'<text x="{x+width+5:.2f}" y="{bar_y+15}">{value:.2f}</text>')
            old, new = values[n]
            reduction = 100 * (old - new) / old if old else 0
            svg.append(f'<text x="{x}" y="{y+68}">{reduction:.1f}% reduction</text>')
    svg.append('</g></svg>')
    path.with_name("comparison.svg").write_text("\n".join(svg))


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("results", type=Path)
    parser.add_argument("--baseline", type=Path)
    args = parser.parse_args()
    render(args.results)
    if args.baseline:
        compare(args.results, args.baseline)
