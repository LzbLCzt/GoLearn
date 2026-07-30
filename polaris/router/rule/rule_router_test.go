package rule

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/config"
)

const (
	testDestNamespace = "Test"
	testDestService   = "lzb_test" // 被调服务

	testSrcNamespace = "Test"
	testSrcService   = "lzb_rulebase_router" // 主调服务

	requestTimes = 100000
)

// ruleRouterYaml SDK 规则路由配置。
// 关键：在 serviceRouter.chain 中只启用 ruleBasedRouter（规则路由），
// 主调通过请求里的 SourceService 传入。
const ruleRouterYaml = `
consumer:
  serviceRouter:
    chain:
      - ruleBasedRouter
`

// instanceStat 记录被调实例信息与命中次数
type instanceStat struct {
	host     string
	port     uint32
	metadata map[string]string
	count    int
}

// TestRuleRouter_CallerLzbRulebaseRouter 验证规则路由：以 lzb_rulebase_router/Test 为主调，
// 请求被调服务 lzb_test/Test，统计最终命中的实例分布（含 metadata，便于确认是否按规则路由）。
func TestRuleRouter_CallerLzbRulebaseRouter(t *testing.T) {
	// 通过内联 yaml 加载 SDK 配置（仅启用规则路由）
	cfg, err := config.LoadConfiguration([]byte(ruleRouterYaml))
	if err != nil {
		t.Fatalf("fail to load polaris config, err: %v", err)
	}

	// 验证本地 polaris-go 的日志能落到默认的 ./polaris/log 目录。
	// 注意：rule.go 中打印 DestRouteRule 的是 Debugf，日志级别必须 <= DebugLog 才能输出。
	_ = api.SetLoggersLevel(api.DebugLog)

	consumer, err := api.NewConsumerAPIByConfig(cfg)
	if err != nil {
		t.Fatalf("fail to create ConsumerAPI, err: %v", err)
	}
	defer consumer.Destroy()

	// 等待服务端配置（路由规则）更新
	time.Sleep(5 * time.Second)

	stats := make(map[string]*instanceStat, 20)
	for i := 0; i < requestTimes; i++ {
		req := &api.GetOneInstanceRequest{}
		req.Namespace = testDestNamespace
		req.Service = testDestService
		// 设置主调信息，规则路由依赖 SourceService 做匹配
		//req.SourceService = &model.ServiceInfo{
		//	Namespace: testSrcNamespace,
		//	Service:   testSrcService,
		//	Metadata: map[string]string{
		//		"k1": "v1",
		//	},
		//}
		rsp, err := consumer.GetOneInstance(req)
		if err != nil {
			t.Fatalf("fail to get one instance at #%d, err: %v", i, err)
		}
		if len(rsp.Instances) == 0 {
			t.Fatalf("empty instances at #%d", i)
		}
		ins := rsp.Instances[0]
		key := fmt.Sprintf("%s:%d", ins.GetHost(), ins.GetPort())
		s, ok := stats[key]
		if !ok {
			s = &instanceStat{
				host:     ins.GetHost(),
				port:     ins.GetPort(),
				metadata: ins.GetMetadata(),
			}
			stats[key] = s
		}
		s.count++
		//time.Sleep(2 * time.Millisecond)
	}

	printStats(t, stats)
}

// printStats 按 host 升序打印命中统计（含 metadata）
func printStats(t *testing.T, stats map[string]*instanceStat) {
	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	total := 0
	t.Logf("=== 请求分配统计（主调：%s/%s，被调：%s/%s，共 %d 次）===",
		testSrcNamespace, testSrcService, testDestNamespace, testDestService, requestTimes)
	for _, k := range keys {
		s := stats[k]
		t.Logf("  %-16s  meta=%v  count=%4d",
			fmt.Sprintf("%s:%d", s.host, s.port), s.metadata, s.count)
		total += s.count
	}
	t.Logf("命中实例总数: %d, 请求总数: %d", len(keys), total)
}
