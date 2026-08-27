#!/bin/bash
set -o pipefail

# 参数：$1 上报次数，$2 上报时间间隔(带单位，如 1s / 500ms)
COUNT=${1:-1}
INTERVAL=${2:-3s}

# 固定业务参数
#SERVER_ADDR="9.134.117.127:8101"
SERVER_ADDR="9.141.110.176:8181"
NAMESPACE="Test"
SERVICE="lzb_healthcheck"
TOKEN="5280e43390604b0d85b8e2f6c592f679"
INSTANCE_ADDR="1.1.1.1:8080"

LOG_FILE="$(dirname "$0")/heartbeat.log"
rm -f "${LOG_FILE}"

echo ">>> run heartbeat grpc client: server=${SERVER_ADDR} service=${SERVICE} count=${COUNT} interval=${INTERVAL}"
go run ./polaris/healthcheck/heartbeat \
    -server="${SERVER_ADDR}" \
    -namespace="${NAMESPACE}" \
    -service="${SERVICE}" \
    -token="${TOKEN}" \
    -addr="${INSTANCE_ADDR}" \
    -count="${COUNT}" \
    -interval="${INTERVAL}" 2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
