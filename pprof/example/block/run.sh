#!/usr/bin/env bash
# 一键跑 block profile 分析
#
# 用法：
#   bash run.sh              # 默认：启动服务 -> 采集 15s block profile -> 打开 Web UI
#   bash run.sh top          # 采集完直接命令行打印 top（不弹浏览器）
#   bash run.sh clean        # 清理生成物
#
# 依赖：go、curl；如果想看调用图，最好装一下 graphviz（brew install graphviz）

set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"

BIN="./blockdemo"
PID_FILE="./blockdemo.pid"
LOG_FILE="./blockdemo.log"
PROFILE_FILE="./block.pprof"
PPROF_URL="http://localhost:6060/debug/pprof/block"
DURATION="${DURATION:-15}"       # 采集时长（秒），block profile 是"当前累计快照"，时间越长累计越多
WEB_PORT="${WEB_PORT:-8080}"     # go tool pprof -http 端口

action="${1:-run}"

stop_prev() {
  if [[ -f "$PID_FILE" ]]; then
    old_pid="$(cat "$PID_FILE" || true)"
    if [[ -n "${old_pid}" ]] && kill -0 "${old_pid}" 2>/dev/null; then
      echo "[run.sh] stopping previous demo pid=${old_pid}"
      kill "${old_pid}" 2>/dev/null || true
      sleep 1
    fi
    rm -f "$PID_FILE"
  fi
  # 兜底：确保 :6060 空闲
  if lsof -iTCP:6060 -sTCP:LISTEN -n -P >/dev/null 2>&1; then
    echo "[run.sh] :6060 still in use, killing occupier"
    lsof -tiTCP:6060 -sTCP:LISTEN | xargs -r kill -9 || true
    sleep 1
  fi
}

case "$action" in
  clean)
    stop_prev
    rm -f "$BIN" "$LOG_FILE" "$PROFILE_FILE"
    echo "[run.sh] cleaned."
    exit 0
    ;;

  run|top)
    stop_prev

    echo "[run.sh] building..."
    go build -o "$BIN" .

    echo "[run.sh] starting demo, log -> $LOG_FILE"
    nohup "$BIN" > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    demo_pid="$(cat "$PID_FILE")"

    # 确保退出时清掉后台进程
    trap 'echo "[run.sh] cleaning up..."; kill '"$demo_pid"' 2>/dev/null || true; rm -f '"$PID_FILE"'' EXIT

    # 等 pprof endpoint ready
    echo "[run.sh] waiting for pprof endpoint ..."
    for i in {1..20}; do
      if curl -fs "http://localhost:6060/debug/pprof/" >/dev/null 2>&1; then
        break
      fi
      sleep 0.3
    done

    # 让 workload 先跑一会儿累计阻塞样本
    echo "[run.sh] warming up ${DURATION}s to accumulate block samples ..."
    sleep "$DURATION"

    echo "[run.sh] fetching block profile -> $PROFILE_FILE"
    curl -fsS "$PPROF_URL" -o "$PROFILE_FILE"
    ls -lh "$PROFILE_FILE"

    if [[ "$action" == "top" ]]; then
#      echo "[run.sh] pprof top (block):"
#      go tool pprof -top -nodecount=15 "$BIN" "$PROFILE_FILE"

      echo
      echo "[run.sh] pprof top -cum (block):"
      go tool pprof -top -cum -nodecount=15 "$BIN" "$PROFILE_FILE"

      echo
      echo "[run.sh] pprof list runWorkload:"
      go tool pprof -list 'runWorkload|consumer|fanInWorker|barrierWorker' "$BIN" "$PROFILE_FILE" || true
    else
      echo "[run.sh] launching web UI at http://localhost:${WEB_PORT} (Ctrl+C to quit)"
      go tool pprof -http=":${WEB_PORT}" "$BIN" "$PROFILE_FILE"
    fi
    ;;

  *)
    echo "unknown action: $action"
    echo "usage: bash run.sh [run|top|clean]"
    exit 1
    ;;
esac
