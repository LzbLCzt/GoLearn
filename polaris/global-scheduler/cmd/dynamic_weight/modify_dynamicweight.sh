#!/bin/bash
# 通过 pb 类型生成 calculate_config，再用 jq 组装外层 body 调用 dynamic-weight 接口
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# 1) 用 Go + protojson 生成 calculate_config（mode 为数字模式）
calculate_config="$(cd "${ROOT_DIR}" && go run ./polaris/global-scheduler/cmd/dynamic_weight/gencalcconfig)"

# 2) scrape_config 字段较少，直接用 jq 内联生成
scrape_config="$(jq -nc '{
  framework:    "vllm",
  metrics_path: "/metrics",
  interval_ms:  5000,
  timeout_ms:   3000,
  metrics_port: 8080
}')"

# 3) 用 jq 组装最外层 body：calculate_config / scrape_config 都以字符串嵌入 params
body="$(jq -nc \
  --arg calc   "${calculate_config}" \
  --arg scrape "${scrape_config}" \
  '[{
    service:   "haixinshi.llm.test",
    namespace: "Test",
    isEnable:  true,
    interval:  10,
    params: {
      metric_source:    "scrape",
      scrape_config:    $scrape,
      calculate_config: $calc
    },
    isUDFEnable:   false,
    language:      "",
    udf:           "",
    service_token: "ec60438a932841d8b32c5a97e20921fd"
  }]')"

echo ">>> request body:"
echo "${body}" | jq .

# 4) 调用 dynamic-weight 接口
curl --location --request PUT 'http://30.163.76.56:8080/naming/v1/dynamicweight' \
  --header 'Platform-Token: a63acf6a46fd44f1ad892f80a2332c13' \
  --header 'Platform-Id: polaris-sdk-test' \
  --header 'Content-Type: application/json' \
  --header 'Cookie: x-client-ssid=bb602415:019ed45e5f72:146ffc; x_host_key_access=97220098170d203108424eb8f2eb97e750352040_s' \
  --data "${body}"
echo ""