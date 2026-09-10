#!/bin/bash
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/bench_dual.log"
INTERVAL=2

# 编译产物放临时目录，脚本退出时自动清理
BIN_DIR="$(mktemp -d -t schedule_client.XXXXXX)"
BIN="${BIN_DIR}/schedule_client"

# 记录当前后台子进程 PID，Ctrl+C 时用于强杀
CHILD_PID=0
STOP=0

cleanup() {
    STOP=1
    # 杀掉正在运行的子进程（包含其子孙进程）
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

echo ">>> 编译 schedule_client -> ${BIN}"
if ! go build -o "${BIN}" ./polaris/global-scheduler/cmd/schedule_client; then
    echo ">>> 编译失败，退出"
    rm -rf "${BIN_DIR}"
    exit 1
fi

echo ">>> 每 ${INTERVAL}s 执行一次 schedule_client，日志: ${LOG_FILE}"
echo ">>> 按 Ctrl+C 退出"

# 单轮执行封装：后台运行 + wait，使 bash 能及时响应 SIGINT
run_once() {
    {
        echo ""
        echo "===== $(date '+%Y-%m-%d %H:%M:%S') 开始执行 ====="
        "${BIN}" -namespace=Test -service=lzb_llm_test 2>&1
        echo "===== $(date '+%Y-%m-%d %H:%M:%S') 执行结束 ====="
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
