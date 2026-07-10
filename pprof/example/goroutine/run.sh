#!/usr/bin/env bash
# 若被 sh(非 bash) 执行，则用 bash 重新执行自己，避免 [[ ]] / {1..N} 等 bashism 语法错误
# 注意：这段用的都是 POSIX 兼容语法，sh 能正确解析这一段
if [ -z "${BASH_VERSION:-}" ]; then
  exec bash "$0" "$@"
fi

# 一键跑 goroutine 泄漏分析
#
# 用法：
#   bash run.sh              # 默认：启动服务 -> 累计 15s -> 打开 Web UI (goroutine profile)
#   bash run.sh top          # 采集完直接命令行打印 top / traces / list（不弹浏览器）
#   bash run.sh diff         # 采两次 profile 做增量对比，最能证明"在泄漏"
#   bash run.sh clean        # 清理生成物
#
# 依赖：go、curl；看调用图建议装 graphviz（brew install graphviz）
#
# 关键 endpoint：
#   /debug/pprof/goroutine            默认 pb.gz 二进制格式，配合 go tool pprof 使用
#   /debug/pprof/goroutine?debug=1    紧凑聚合文本，按栈相同归并计数
#   /debug/pprof/goroutine?debug=2    完整栈（每个 goroutine 一段），最详细

set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"

BIN="./goroutinedemo"
PID_FILE="./goroutinedemo.pid"
LOG_FILE="./goroutinedemo.log"
PROFILE_FILE="./goroutine.pprof"
PROFILE_FILE_1="./goroutine.t1.pprof"
PROFILE_FILE_2="./goroutine.t2.pprof"
STACK_TXT="./goroutine.stack.txt"
PPROF_URL="http://localhost:6060/debug/pprof/goroutine"
STACK_URL="http://localhost:6060/debug/pprof/goroutine?debug=2"
DURATION="${DURATION:-15}"       # 累计时长（秒）
DIFF_GAP="${DIFF_GAP:-10}"       # diff 模式：两次采样之间等多久
WEB_PORT="${WEB_PORT:-8080}"

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
  if lsof -iTCP:6060 -sTCP:LISTEN -n -P >/dev/null 2>&1; then
    echo "[run.sh] :6060 still in use, killing occupier"
    lsof -tiTCP:6060 -sTCP:LISTEN | xargs -r kill -9 || true
    sleep 1
  fi
}

start_demo() {
  echo "[run.sh] building..."
  go build -o "$BIN" .

  echo "[run.sh] starting demo, log -> $LOG_FILE"
  nohup "$BIN" > "$LOG_FILE" 2>&1 &
  echo $! > "$PID_FILE"
  demo_pid="$(cat "$PID_FILE")"

  trap 'echo "[run.sh] cleaning up..."; kill '"$demo_pid"' 2>/dev/null || true; rm -f '"$PID_FILE"'' EXIT

  echo "[run.sh] waiting for pprof endpoint ..."
  for i in {1..20}; do
    if curl -fs "http://localhost:6060/debug/pprof/" >/dev/null 2>&1; then
      break
    fi
    sleep 0.3
  done
}

case "$action" in
  clean)
    stop_prev
    rm -f "$BIN" "$LOG_FILE" "$PROFILE_FILE" "$PROFILE_FILE_1" "$PROFILE_FILE_2" "$STACK_TXT" "./result.txt"
    echo "[run.sh] cleaned."
    exit 0
    ;;

  run|top)
    stop_prev
    start_demo

    echo "[run.sh] warming up ${DURATION}s to accumulate leaked goroutines ..."
    sleep "$DURATION"

    echo "[run.sh] fetching goroutine profile -> $PROFILE_FILE"
    curl -fsS "$PPROF_URL" -o "$PROFILE_FILE"
    echo "[run.sh] fetching full stacks (debug=2) -> $STACK_TXT"
    curl -fsS "$STACK_URL" -o "$STACK_TXT"
    ls -lh "$PROFILE_FILE" "$STACK_TXT"

    if [[ "$action" == "top" ]]; then
      RESULT_FILE="./result.txt"
      rm -f "$RESULT_FILE"                          # 每次先清理旧结果
      echo "[run.sh] result will be saved to $RESULT_FILE"

      {
        echo
        echo "[run.sh] pprof top -cum (goroutine):"
        go tool pprof -top -cum -nodecount=15 "$BIN" "$PROFILE_FILE"

        echo
        echo "[run.sh] pprof traces (each unique goroutine stack + count):"
        go tool pprof -traces "$BIN" "$PROFILE_FILE"

        echo
        echo "[run.sh] pprof list leak functions:"
        go tool pprof -list 'leakOnUnbufferedSend|leakOnRecvNoClose|leakOnForgottenTimer' "$BIN" "$PROFILE_FILE" || true

        echo
        echo "[run.sh] 泄漏 goroutine 分组统计（来自 debug=2 全栈）："
        # 抓每段 goroutine 的第一行 "goroutine N [state, ...]"，
        # 再拼上栈里第一个 main.* 的函数名作为分组键
        awk '
          /^goroutine [0-9]+ \[/ {
            state=$0; sub(/^goroutine [0-9]+ \[/, "", state); sub(/\].*/, "", state);
            fn=""; next
          }
          fn=="" && /^main\./ {
            fn=$1; sub(/\(.*/, "", fn);
            print state " | " fn
          }
        ' "$STACK_TXT" | sort | uniq -c | sort -rn
      } 2>&1 | tee "$RESULT_FILE"
    else
      echo "[run.sh] launching web UI at http://localhost:${WEB_PORT} (Ctrl+C to quit)"
      go tool pprof -http=":${WEB_PORT}" "$BIN" "$PROFILE_FILE"
    fi
    ;;

  diff)
    stop_prev
    start_demo

    echo "[run.sh] [t1] first snapshot after 3s ..."
    sleep 3
    curl -fsS "$PPROF_URL" -o "$PROFILE_FILE_1"

    echo "[run.sh] waiting ${DIFF_GAP}s to let leakers grow ..."
    sleep "$DIFF_GAP"

    echo "[run.sh] [t2] second snapshot ..."
    curl -fsS "$PPROF_URL" -o "$PROFILE_FILE_2"

    ls -lh "$PROFILE_FILE_1" "$PROFILE_FILE_2"

    echo
    echo "[run.sh] pprof diff (t2 - t1)  →  正值 = 期间新增的泄漏 goroutine："
    go tool pprof -top -cum -nodecount=15 -diff_base "$PROFILE_FILE_1" "$BIN" "$PROFILE_FILE_2"
    ;;

  *)
    echo "unknown action: $action"
    echo "usage: bash run.sh [run|top|diff|clean]"
    exit 1
    ;;
esac
