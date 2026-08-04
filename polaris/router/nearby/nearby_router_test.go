package nearby

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/config"
)

const (
	testNamespace = "Test"
	testService   = "lzb_nearby_router"

	callerRegion = "华南"
	callerZone   = "深圳"
	callerCampus = "深圳一区"

	requestTimes = 10000
)

// nearbyRouterYaml SDK 就近路由配置。
// 关键：polaris-go 的 unhealthyPercentToDegrade 只读本地配置（默认 100，即 100% 才降级），
// 服务端控制台配置不生效，必须在这里显式覆盖，才能与控制台配置对齐。
const nearbyRouterYaml = `
consumer:
  serviceRouter:
    plugin:
      nearbyBasedRouter:
        matchLevel: zone                       # 就近级别：城市
        maxMatchLevel: ""                      # 最大降级级别：允许降到全部（all 用空串表示）
        enableDegradeByUnhealthyPercent: true  # 开启按不健康比例降级
        unhealthyPercentToDegrade: 40          # 不健康比例 >= 40% 触发降级
`

// instanceStat 记录被调实例的 CMDB 信息与命中次数
type instanceStat struct {
	host   string
	port   uint32
	region string
	zone   string
	campus string
	count  int
}

// TestNearbyRouter_AllHealthy_CallerInShenzhen
func TestNearbyRouter_AllHealthy_CallerInShenzhen(t *testing.T) {
	// 通过内联 yaml 加载 SDK 配置（覆盖 nearby 路由默认值）
	cfg, err := config.LoadConfiguration([]byte(nearbyRouterYaml))
	if err != nil {
		t.Fatalf("fail to load polaris config, err: %v", err)
	}

	// 代码方式设置主调位置信息
	loc := cfg.GetGlobal().GetAPI().GetLocation()
	loc.SetEnableUpdate(false) // 关闭定时更新自身地域信息
	loc.SetRegion(callerRegion)
	loc.SetZone(callerZone)
	loc.SetCampus(callerCampus)

	consumer, err := api.NewConsumerAPIByConfig(cfg)
	if err != nil {
		t.Fatalf("fail to create ConsumerAPI, err: %v", err)
	}
	defer consumer.Destroy()

	// 等待服务端配置更新
	time.Sleep(5 * time.Second)

	stats := make(map[string]*instanceStat, 20)
	for i := 0; i < requestTimes; i++ {
		req := &api.GetOneInstanceRequest{}
		req.Namespace = testNamespace
		req.Service = testService
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
				host:   ins.GetHost(),
				port:   ins.GetPort(),
				region: ins.GetRegion(),
				zone:   ins.GetZone(),
				campus: ins.GetCampus(),
			}
			stats[key] = s
		}
		s.count++
		time.Sleep(2 * time.Millisecond)
	}

	printStats(t, stats)
}

// printStats 按 host 升序打印命中统计（含 CMDB 信息）
func printStats(t *testing.T, stats map[string]*instanceStat) {
	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	total := 0
	t.Logf("=== 请求分配统计（主调：%s/%s/%s，共 %d 次）===",
		callerRegion, callerZone, callerCampus, requestTimes)
	t.Logf("  %-16s  %-8s  %-8s  %-12s  %s", "HOST:PORT", "REGION", "ZONE", "CAMPUS", "COUNT")
	for _, k := range keys {
		s := stats[k]
		t.Logf("  %-16s  %-8s  %-8s  %-12s  %4d",
			fmt.Sprintf("%s:%d", s.host, s.port), s.region, s.zone, s.campus, s.count)
		total += s.count
	}
	t.Logf("命中实例总数: %d, 请求总数: %d", len(keys), total)
}
