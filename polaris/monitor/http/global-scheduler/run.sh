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

# 通过 HTTP 接口 /monitor/v1/CollectScheduleRequestStat 上报全局调度请求统计
# 对应服务端：server/polaris/http/server.go -> Server.CollectScheduleRequestStat
set -o pipefail

# HTTP 服务地址（端口请按实际部署配置修改，默认与 gRPC 同 IP）
HTTP_ADDR="${HTTP_ADDR:-9.134.117.127:8087}"

curl -s -X POST "http://${HTTP_ADDR}/monitor/v1/CollectScheduleRequestStat" \
    -H 'Content-Type: application/json' \
    -d '{
  "id": "test-schedule-stat-1",
  "sdk_token": {
    "ip": "127.0.0.1",
    "client": "polaris-go",
    "version": "2.6.0"
  },
  "service": "lzb_test3",
  "namespace": "Test",
  "load_balance_type": "GLOBAL_P2C",
  "request_stat": [
    {"stat_type": "RequestSuccess", "period_times": 1, "reason": "ok"},
    {"stat_type": "RequestDegrade", "period_times": 1, "reason": "degrade"},
    {"stat_type": "RequestFailed", "period_times": 1, "reason": "failed"}
  ],
  "time": "2026-08-04T12:00:00Z"
}'

echo ""

