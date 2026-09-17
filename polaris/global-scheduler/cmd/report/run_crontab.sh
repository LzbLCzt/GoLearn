#!/bin/bash
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/report.log"
INTERVAL=60

# 编译产物放临时目录，脚本退出时自动清理
BIN_DIR="$(mktemp -d -t report.XXXXXX)"
BIN="${BIN_DIR}/report"

# 记录当前后台子进程 PID，Ctrl+C 时用于强杀
CHILD_PID=0
STOP=0

cleanup() {
    STOP=1
    if [[ "${CHILD_PID}" -gt 0 ]] && kill -0 "${CHILD_PID}" 2>/dev/null; then
        kill -TERM -"${CHILD_PID}" 2>/dev/null || kill -TERM "${CHILD_PID}" 2>/dev/null
        sleep 0.2
        kill -KILL -"${CHILD_PID}" 2>/dev/null || kill -KILL "${CHILD_PID}" 2>/dev/null
    fi
    rm -rf "${BIN_DIR}"
    echo ""
    echo ">>> 已停止定时任务，完整日志: ${LOG_FILE}"
    exit 0
}
trap cleanup INT TERM

# 首次启动清空日志，如需保留历史请注释下一行
: > "${LOG_FILE}"

echo ">>> 编译 report -> ${BIN}"
if ! go build -o "${BIN}" ./polaris/global-scheduler/cmd/Report; then
    echo ">>> 编译失败，退出"
    rm -rf "${BIN_DIR}"
    exit 1
fi

echo ">>> 每 ${INTERVAL}s 上报一次，日志: ${LOG_FILE}"
echo ">>> 按 Ctrl+C 退出"

# 单轮执行封装：后台运行 + wait，使 bash 能及时响应 SIGINT
run_once() {
    {
        echo ""
        echo "===== $(date '+%Y-%m-%d %H:%M:%S') 开始上报 ====="
        "${BIN}" \
            -namespace=Test \
            -service=lzb_llm_test \
            -instance="host=1.1.1.1,port=8080,kv_cache_usage_perc=0.3,num_requests_running=50,num_requests_waiting=100" \
            -instance="host=1.1.1.2,port=8080,kv_cache_usage_perc=0.7,num_requests_running=100,num_requests_waiting=70" \
            -instance="host=1.1.1.3,port=8080,kv_cache_usage_perc=0.6,num_requests_running=200,num_requests_waiting=1000" \
            2>&1
        echo "===== $(date '+%Y-%m-%d %H:%M:%S') 上报结束 ====="
    } | tee -a "${LOG_FILE}" &
    CHILD_PID=$!
    wait "${CHILD_PID}" 2>/dev/null
    CHILD_PID=0
}

while [[ "${STOP}" -eq 0 ]]; do
    run_once
    [[ "${STOP}" -eq 1 ]] && break

    # sleep 也放后台 + wait，Ctrl+C 立即中断
    sleep "${INTERVAL}" &
    CHILD_PID=$!
    wait "${CHILD_PID}" 2>/dev/null
    CHILD_PID=0
done
