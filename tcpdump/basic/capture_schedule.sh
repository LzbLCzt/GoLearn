#!/bin/bash
# 抓取 CollectScheduleRequestStat gRPC 接口流量的 tcpdump 脚本
# 接口协议：gRPC over HTTP/2，服务端地址 9.134.117.127:8086
# 用法：bash capture_schedule.sh [网卡名]   （不传则默认抓所有网卡 any）
set -o pipefail

INTERFACE="${1:-any}"                 # 网卡：默认 any（所有网卡）；本机回环可用 lo
HOST="9.134.117.127"                  # gRPC 服务端 IP
PORT="8086"                           # gRPC 监听端口
PCAP_FILE="schedule_capture.pcap"     # 保存的 pcap 文件，可用 wireshark/tshark 分析

echo ">>> 目标过滤: host ${HOST} and port ${PORT}"
echo ">>> 网卡=${INTERFACE}  保存=${PCAP_FILE}"
echo ">>> 按 Ctrl+C 提前结束；结束后可用: tshark -r ${PCAP_FILE} -Y http2"

# -i 网卡  -nn 不解析主机名/端口名（显示纯 IP:端口）
# -s0 抓完整包（默认只抓前 68 字节，会截断 HTTP/2 帧）
# -A  以 ASCII 打印包内容（方便直接看到 HTTP/2 头部明文）
# -w  同时写 pcap 文件，供后续用图形化工具深入分析
tcpdump -i "${INTERFACE}" -nn -s0 -A \
    "host ${HOST} and port ${PORT}" -w "${PCAP_FILE}"

echo ""
echo ">>> 抓包结束，pcap 已保存: $(pwd)/${PCAP_FILE}"
