#!/bin/bash
set -o pipefail

LOG_FILE="$(dirname "$0")/report.log"
rm -f "${LOG_FILE}"

go run ./polaris/global-scheduler/cmd/Report \
    -namespace=Test \
    -service=haixinshi.llm.test \
    -instance="host=21.91.241.157,port=8000" \
    -instance="host=11.12.13.14,port=8080" \
    -instance="host=30.163.16.111,port=8000" \
    -instance="host=11.12.13.15,port=8080" \
    -instance="host=11.12.13.16,port=8080" \
    -instance="host=11.12.13.17,port=8080" \
    2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"

# chmod +x polaris/global-scheduler/cmd/Report/run.sh && bash polaris/global-scheduler/cmd/Report/run.sh