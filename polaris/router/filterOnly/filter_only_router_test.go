package filterOnly

import (
	"testing"
	"time"

	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/config"
)

const (
	testNamespace = "Test"
	testService   = "lzb_filter_only_router" // 被调服务（后置路由插件测试）
)

// filterOnlyRouterYaml 仅保留默认路由链，filterOnlyRouter 后置过滤由 SDK 自动执行
const filterOnlyRouterYaml = `
consumer:
  serviceRouter:
    chain:
      - ruleBasedRouter
      - nearbyBasedRouter
`

// TestFilterOnlyRouter 验证后置路由插件：请求被调服务 lzb_filter_only_router/Test，
// 触发 SDK 的后置 filterOnlyRouter 过滤，并打印最终命中的实例信息。
func TestFilterOnlyRouter(t *testing.T) {
	cfg, err := config.LoadConfiguration([]byte(filterOnlyRouterYaml))
	if err != nil {
		t.Fatalf("fail to load polaris config, err: %v", err)
	}

	consumer, err := api.NewConsumerAPIByConfig(cfg)
	if err != nil {
		t.Fatalf("fail to create ConsumerAPI, err: %v", err)
	}
	defer consumer.Destroy()

	// 等待服务端配置（路由规则）下发
	time.Sleep(5 * time.Second)

	req := &api.GetOneInstanceRequest{}
	req.Namespace = testNamespace
	req.Service = testService
	rsp, err := consumer.GetOneInstance(req)
	if err != nil {
		t.Fatalf("fail to get one instance, err: %v", err)
	}
	if len(rsp.Instances) == 0 {
		t.Fatalf("empty instances")
	}

	ins := rsp.Instances[0]
	t.Logf("=== 后置路由插件测试结果（被调：%s/%s）===", testNamespace, testService)
	t.Logf("命中实例: host=%s, port=%d, metadata=%v, region=%s, zone=%s, campus=%s",
		ins.GetHost(), ins.GetPort(), ins.GetMetadata(),
		ins.GetRegion(), ins.GetZone(), ins.GetCampus())
}
