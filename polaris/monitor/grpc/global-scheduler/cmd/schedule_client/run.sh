#!/bin/bash
set -o pipefail

LOG_FILE="$(dirname "$0")/schedule_client.log"
rm -f "${LOG_FILE}"
go run ./polaris/monitor/grpc/global-scheduler/cmd/schedule_client \
    -addr=9.134.117.127:8086 \
    -namespace=Test -service=lzb_test2 -lb=GLOBAL_P2C \
    -count=3 -dial-timeout=15s -call-timeout=15s 2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
