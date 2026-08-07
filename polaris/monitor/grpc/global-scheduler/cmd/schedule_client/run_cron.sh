#!/bin/bash
#
# Tencent is pleased to support the open source community by making Polaris available.
#
# Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
#
# Licensed under the BSD 3-Clause License (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# https://opensource.org/licenses/BSD-3-Clause
#
# Unless required by applicable law or agreed to in writing, software distributed
# under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.

# 定时每分钟调用一次 schedule_client，并在每次请求前后打印当前时间
set -o pipefail

# 切换到仓库根目录，保证 go run 的相对路径正确
cd "$(git rev-parse --show-toplevel)" || exit 1

LOG_FILE="$(dirname "$0")/schedule_client_cron.log"
rm -f "${LOG_FILE}"

while true; do
   echo "[$(date '+%Y-%m-%d %H:%M:%S')] >>> 开始调用 schedule_client"
      go run ./polaris/monitor/grpc/global-scheduler/cmd/schedule_client \
          -addr=11.140.163.248:8089 \
          -namespace=Test -service=lzb_test2 -lb=GLOBAL_P2C \
          -count=58 -dial-timeout=15s -call-timeout=15s 2>&1 | tee -a "${LOG_FILE}"
      echo "[$(date '+%Y-%m-%d %H:%M:%S')] <<< 本次调用完成，60 秒后重试"
      sleep 60

    echo "[$(date '+%Y-%m-%d %H:%M:%S')] >>> 开始调用 schedule_client"
    go run ./polaris/monitor/grpc/global-scheduler/cmd/schedule_client \
        -addr=9.146.200.123:8089 \
        -namespace=Test -service=lzb_test2 -lb=GLOBAL_P2C \
        -count=58 -dial-timeout=15s -call-timeout=15s 2>&1 | tee -a "${LOG_FILE}"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] <<< 本次调用完成，60 秒后重试"
    sleep 60
done

