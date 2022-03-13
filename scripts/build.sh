#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if ! [[ "$0" =~ scripts/build.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

if [[ $# -eq 1 ]]; then
    binary_path=$1
else
    echo "Invalid arguments to build subnet_evm. Requires one arguments to specify binary location."
    exit 1
fi

# Build spacesvm, which is run as a subprocess
mkdir -p ./build

echo "Building spacesvm in $binary_path"
go build -o "$binary_path" ./cmd/spacesvm

echo "Building spaces-cli in ./build/spaces-cli"
go build -o ./build/spaces-cli ./cmd/spaces-cli
