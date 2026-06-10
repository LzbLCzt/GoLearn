#!/bin/bash
set -o pipefail

LOG_FILE="$(dirname "$0")/bench_dual.log"
rm "${LOG_FILE}"
go run ./polaris/global-scheduler/cmd/schedule_bench_dual \
    -namespace=Test -service=lzb_test 2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
