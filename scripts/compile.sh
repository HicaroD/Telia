#!/bin/bash

set -e          # Exits immediately when some command has a non-zero status
set -x          # Debug mode
set -u          # Error when referencing undefined variable
set -o pipefail # Fail when any error ocurrs on pipe operations

TELIA_FILE=$1
LLVM_FILE=telia.ll

echo "Compiling program"
go build -tags=llvm18
./Telia $TELIA_FILE
echo "Generating binary executable"
clang -O3 -Wall $LLVM_FILE
./a.out
