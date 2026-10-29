#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/.."
run="bench/results/$(date -u +%Y%m%dT%H%M%SZ)-$$"
mkdir -p "$run"
go build -o "$run/profile" ./bench
"$run/profile" -out "$run" "$@"
go tool pprof -top "$run/profile" "$run/cpu.pprof" > "$run/cpu.txt"
go tool pprof -top -alloc_space "$run/profile" "$run/allocs.pprof" > "$run/allocs.txt"
python3 bench/charts.py "$run/results.json"
echo "Results: $run"
