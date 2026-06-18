#!/bin/bash
set -o pipefail

LOG_FILE="$(dirname "$0")/query_service_load.log"
rm -f "${LOG_FILE}"
go run ./polaris/global-scheduler/cmd/query_service_load \
    -namespace=Test -service=haixinshi.llm.test 2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
