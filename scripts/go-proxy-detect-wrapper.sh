#!/usr/bin/env bash
set -euo pipefail

IMAGE="akamaiopen-edgegrid-golang/base"

exec docker run --rm \
  -v "$PWD:$PWD" \
  -w "$PWD" \
  "$IMAGE" \
  /usr/local/bin/go-proxy "$@"