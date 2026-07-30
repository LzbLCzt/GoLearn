#!/bin/bash
set -o pipefail

LOG_FILE="$(dirname "$0")/report.log"
rm -f "${LOG_FILE}"

go run ./polaris/global-scheduler/cmd/Report \
    -namespace=Test \
    -service=haixinshi.llm.test \
    -instance="host=21.91.241.157,port=8000,kv_cache_usage_perc=0.7,num_requests_running=100,num_requests_waiting=100" \
    -instance="host=11.12.13.14,port=8080,kv_cache_usage_perc=0.4,num_requests_running=20,num_requests_waiting=80" \
    -instance="host=30.163.16.111,port=8000,kv_cache_usage_perc=0.3,num_requests_running=30,num_requests_waiting=300" \
    -instance="host=11.12.13.15,port=8080,kv_cache_usage_perc=0.2,num_requests_running=60,num_requests_waiting=40" \
    -instance="host=11.12.13.16,port=8080,kv_cache_usage_perc=0.5,num_requests_running=50,num_requests_waiting=200" \
#    -instance="host=11.12.13.17,port=8080,kv_cache_usage_perc=0.3,num_requests_running=20,num_requests_waiting=0" \
    2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
