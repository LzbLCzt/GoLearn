#!/bin/bash
set -o pipefail

# 依赖的 polaris-go 项目路径（GoLearn 的 go.mod replace 指向此）
POLARIS_GO_DIR="/Users/zhengbangli/code/polaris-go"

# 每次运行前，先 build 一下依赖的 polaris-go 的 api 包（业务实际引用的入口）
# 备注：polaris-go 下老版本 gomonkey 在 arm64/新 Go 上 build 会失败，
#       仅业务方使用的 api 子包不受影响，因此只 build ./api/...
echo ">>> build polaris-go: ${POLARIS_GO_DIR}/api/..."
(cd "${POLARIS_GO_DIR}" && go build ./api/...) || {
    echo ">>> build polaris-go failed"
    exit 1
}

# select 次数（可通过第一个参数覆盖）
COUNT=${1:-20}

LOG_FILE="$(dirname "$0")/basic.log"
rm -f "${LOG_FILE}"

echo ">>> run circuitbreaker basic, count=${COUNT}"
go run ./polaris/circuitbreaker/basic \
    -namespace=Test -service=lzb_test \
    -src-namespace=Test -src-service=lzb_test2 \
    -count="${COUNT}" -interval=1s 2>&1 | tee "${LOG_FILE}"

echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
