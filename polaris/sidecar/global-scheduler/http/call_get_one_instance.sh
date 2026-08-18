#!/bin/bash

# 每秒调用一次 GetOneInstance 接口

URL="http://9.134.117.127:8094/v1/GetOneInstance"
DATA='{"namespace":"Test","service":"lzb_llm_half_minute_test"}'

while true; do
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] 调用接口..."
    curl -s -X POST "${URL}" \
        -H "Content-Type: application/json" \
        -d "${DATA}"
    echo ""
    sleep 1
done
