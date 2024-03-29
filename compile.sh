#!/bin/bash

set -e          # Exits immediately when some command has a non-zero status
set -x          # Debug mode
set -u          # Error when referencing undefined variable
set -o pipefail # Fail when any error ocurrs on pipe operations

TELIA_FILE=$1
LLVM_FILE=telia.ll
LLVM_ASSEMBLY_FILE=telia.s

echo "Compiling program"
go build -tags=llvm16
./telia-lang $TELIA_FILE
echo "Generating binary executable"
clang -fomit-frame-pointer $LLVM_FILE
./a.out
