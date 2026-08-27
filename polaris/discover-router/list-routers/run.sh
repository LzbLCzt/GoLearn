#!/bin/bash
set -o pipefail

# discover-router 服务地址（可按需修改）
SERVER_ADDR="9.134.117.127:8098"

# 请求参数
NAMESPACE="Test"
SERVICE="lzb_test"
# ServiceRouterReq.ClusterType：服务路由使用的集群类型
# 0=Unknown 1=Discover 2=Heartbeat 3=OpenApi 4=Register 5=Domain 6=Report
SERVICE_CLUSTER_TYPE=1
# ResourceRouterReq.ClusterType：资源路由使用的集群类型
RESOURCE_CLUSTER_TYPE=3
# ResourceType: 0=Unknown 1=InstanceId 2=RateLimitId 3=DomainId
RESOURCE_TYPE=2
# ResourceId: 资源ID，查询资源路由时需传入
#RESOURCE_ID="050eeee6f0628f806d59cca043c658502b9b8028"  # 实例id
#RESOURCE_ID="68240b3a236a4d56a2c7a75d831f400e"  # 限流规则id
RESOURCE_ID="43033c765a294032bf76508ed1cf25a4"  # 域名id

LOG_FILE="$(dirname "$0")/list_routers_demo.log"
rm -f "${LOG_FILE}"

# 查询服务路由
#go run ./polaris/discover-router/list-routers \
#    -addr="${SERVER_ADDR}" \
#    -namespace="${NAMESPACE}" \
#    -service="${SERVICE}" \
#    -service_cluster_type="${SERVICE_CLUSTER_TYPE}"  2>&1 | tee "${LOG_FILE}"

# 查询资源路由
go run ./polaris/discover-router/list-routers \
    -addr="${SERVER_ADDR}" \
    -namespace="${NAMESPACE}" \
    -service="${SERVICE}" \
    -service_cluster_type="${SERVICE_CLUSTER_TYPE}" \
    -resource_cluster_type="${RESOURCE_CLUSTER_TYPE}" \
    -resource_type="${RESOURCE_TYPE}" \
    -resource_id="${RESOURCE_ID}" 2>&1 | tee "${LOG_FILE}"
echo ""
echo ">>> 完整日志已保存: ${LOG_FILE}"
