#!/bin/bash
set -e

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
cd "${SCRIPT_DIR}"

go run . \
  -addr=9.134.117.127:8092 \
  -namespace=Test \
  -service=lzb_test2
