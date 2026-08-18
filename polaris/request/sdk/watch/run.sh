#!/bin/bash
set -o pipefail

LOG_FILE="$(dirname "$0")/watch.log"
rm -f "${LOG_FILE}"

go run ./polaris/request/sdk/watch \
    -namespace=Test \
    -service=lzb_test \
    2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
