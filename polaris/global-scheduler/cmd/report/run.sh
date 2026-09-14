#!/bin/bash
set -o pipefail

LOG_FILE="$(dirname "$0")/report.log"
rm -f "${LOG_FILE}"

# `-a`：强制重新编译所有依赖包，**禁用增量缓存复用**，全部从头编译
go run -a ./polaris/global-scheduler/cmd/Report \
    -namespace=Test \
    -service=lzb_llm_test \
    -instance="host=1.1.1.1,port=8080,kv_cache_usage_perc=0.32,num_requests_running=50,num_requests_waiting=100" \
    -instance="host=1.1.1.2,port=8080,kv_cache_usage_perc=0.7,num_requests_running=100,num_requests_waiting=70" \
    -instance="host=1.1.1.3,port=8080,kv_cache_usage_perc=0.6,num_requests_running=20,num_requests_waiting=100" \
    2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
