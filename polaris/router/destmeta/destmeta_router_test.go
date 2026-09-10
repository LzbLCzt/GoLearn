package destmeta

import (
	"fmt"
	"sort"
	"testing"

	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/config"
	"git.woa.com/polaris/polaris-go/v2/pkg/model"
)

const (
	testDestNamespace = "Test"
	testDestService   = "lzb_test" // 被调服务

	requestTimes = 1000
)

// dstMetaRouterYaml SDK 元数据路由配置。
// 关键：serviceRouter.chain 中只启用 dstMetaRouter（元数据路由），
// 被调实例的匹配条件通过请求里的 Metadata 传入。
// 默认 failover 类型为 none，即无匹配实例时直接报错，最能体现元数据路由的过滤效果。
const dstMetaRouterYaml = `
consumer:
  serviceRouter:
    chain:
      - dstMetaRouter
`

// instanceStat 记录被调实例信息与命中次数
type instanceStat struct {
	host     string
	port     uint32
	metadata map[string]string
	count    int
}

// TestDstMetaRouter 验证元数据路由：请求被调服务 lzb_test/Test，
// 只允许命中 metadata 中包含 k4=v4 的被调实例，统计命中分布。
func TestDstMetaRouter(t *testing.T) {
	cfg, err := config.LoadConfiguration([]byte(dstMetaRouterYaml))
	if err != nil {
		t.Fatalf("fail to load polaris config, err: %v", err)
	}

	// 打开 debug 日志，便于观察 dstMetaRouter 的过滤过程
	_ = api.SetLoggersLevel(api.DebugLog)

	consumer, err := api.NewConsumerAPIByConfig(cfg)
	if err != nil {
		t.Fatalf("fail to create ConsumerAPI, err: %v", err)
	}
	defer consumer.Destroy()

	stats := make(map[string]*instanceStat, 20)
	for i := 0; i < requestTimes; i++ {
		getOneInstanceReq := &api.GetOneInstanceRequest{}
		getOneInstanceReq.Namespace = testDestNamespace
		getOneInstanceReq.Service = testDestService
		// 元数据路由：SDK 只会返回 metadata 中完全包含以下 KV 的被调实例
		dstMetadata := map[string]string{
			"k4": "v4",
		}

		// 元数据路由兜底策略
		// 元数据路由failover策略：GetOneHealth
		//if len(dstMetadata) > 0 {
		//	getOneInstanceReq.Metadata = dstMetadata
		//	getOneInstanceReq.EnableFailOverDefaultMeta = true
		//	getOneInstanceReq.FailOverDefaultMeta = model.FailOverDefaultMetaConfig{
		//		Type: model.GetOneHealth,
		//	}
		//}

		// 元数据路由failover策略：NotContainMetaKey
		if len(dstMetadata) > 0 {
			getOneInstanceReq.Metadata = dstMetadata
			getOneInstanceReq.EnableFailOverDefaultMeta = true
			getOneInstanceReq.FailOverDefaultMeta = model.FailOverDefaultMetaConfig{
				Type: model.NotContainMetaKey,
			}
		}

		// 元数据路由failover策略：CustomMeta
		//if len(dstMetadata) > 0 {
		//	getOneInstanceReq.Metadata = dstMetadata
		//	getOneInstanceReq.EnableFailOverDefaultMeta = true
		//	getOneInstanceReq.FailOverDefaultMeta = model.FailOverDefaultMetaConfig{
		//		Type: model.CustomMeta,
		//		Meta: map[string]string{
		//			"k1": "v3",
		//		},
		//	}
		//}

		rsp, err := consumer.GetOneInstance(getOneInstanceReq)
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
	t.Logf("=== 请求分配统计（被调：%s/%s，匹配 metadata: k4=v4，共 %d 次）===",
		testDestNamespace, testDestService, requestTimes)
	for _, k := range keys {
		s := stats[k]
		t.Logf("  %-16s  meta=%v  count=%4d",
			fmt.Sprintf("%s:%d", s.host, s.port), s.metadata, s.count)
		total += s.count
	}
	t.Logf("命中实例总数: %d, 请求总数: %d", len(keys), total)
}
