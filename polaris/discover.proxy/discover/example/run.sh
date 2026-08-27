#!/bin/bash
set -o pipefail

# discover-proxy 服务器地址（可按需修改）
SERVER_ADDR="9.134.117.127:8089"
NAMESPACE="Test"
#SERVICE="polaris.discover.test"
SERVICE="lzb_test"

LOG_FILE="$(dirname "$0")/discover_demo.log"
rm -f "${LOG_FILE}"
go run ./polaris/discover.proxy/discover/example \
    -addr="${SERVER_ADDR}" -namespace="${NAMESPACE}" -service="${SERVICE}" 2>&1 | tee "${LOG_FILE}"
echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
