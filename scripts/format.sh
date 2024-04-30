#!/bin/bash

set -e          # Exits immediately when some command has a non-zero status
set -x          # Debug mode
set -u          # Error when referencing undefined variable
set -o pipefail # Fail when any error ocurrs on pipe operations

go fmt ./...
golines --write-output .
