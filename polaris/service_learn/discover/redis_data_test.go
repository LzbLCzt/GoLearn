package discover

import (
	"context"
	"log"
	"testing"
	"time"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	v2 "git.woa.com/polaris/polaris-server-api/api/v2/common"
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
)

func InitRedisClient(t *testing.T) *redis.Client {

	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "9.141.69.182:6379",
		Password: "@d@5y2@@polarisM@42021",
		DB:       0,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("Redis connection failed: %+v", err)
	}

	log.Printf("Redis client initialized successfully")
	return rdb
}

func TestRedisData(t *testing.T) {
	// pbSourceStr 从/Users/zhengbangli/code/GoLearn/polaris/service_learn/discover/redis_data_test.go 读取
	pbSourceStr := ""
	svcIns := &v2.DiscoverResponse{}

	if err := proto.Unmarshal([]byte(pbSourceStr), svcIns); err != nil {
		t.Fatal(err)
	}

	t.Log(svcIns)
}

func TestRedisData1(t *testing.T) {
	ctx := context.Background()
	rdb := InitRedisClient(t)
	pbSourceStr, err := rdb.HGet(ctx, "circuit_v2.Test.migration-service-2026-01-08T103258.180797", "info").Bytes()
	if err != nil {
		t.Fatal(err)
	}
	//pbSourceStr := "\n$\n\"test-name-1767839604946582744.Test\x12\x10\n\x0etest-version-0\x1a\x1f\n\x1dtest-name-1767839604946582744\"\x06\n\x04Test*,\n*migration-service-2026-01-08T103258.1807972\x06\n\x04Test:P\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1:P\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1BP\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1BP\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1R\x00Z\x00b\x00j\x00r\x15\n\x132026-01-08 10:33:31z\x15\n\x132026-01-08 10:33:31\x82\x01\"\n 980e46a59db646d884d70dfaf0e99f38"
	circuitBreaker := &apiV1Model.CircuitBreakerV2{}

	if err := proto.Unmarshal([]byte(pbSourceStr), circuitBreaker); err != nil {
		t.Fatal(err)
	}

	t.Log(circuitBreaker)

}

func TestRedisData2(t *testing.T) {
	pbSourceStr := "\n$\n\"test-name-1767839604946582744.Test\x12\x10\n\x0etest-version-0\x1a\x1f\n\x1dtest-name-1767839604946582744\"\x06\n\x04Test*,\n*migration-service-2026-01-08T103258.1807972\x06\n\x04Test:P\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1:P\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1BP\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1BP\x12&\n\x12\n\x10namespace-test-0\x12\x10\n\x0eservice-test-0\x12&\n\x12\n\x10namespace-test-1\x12\x10\n\x0eservice-test-1R\x0e\n\x0clegacy-ownerZ\x11\n\x0flegacy-businessb\x13\n\x11legacy-departmentj\x10\n\x0elegacy-commentr\x15\n\x132026-01-07 10:33:31z\x15\n\x132026-01-08 10:33:31\x82\x01\"\n 980e46a59db646d884d70dfaf0e99f38"
	circuitBreaker := &apiV1Model.CircuitBreakerV2{}

	if err := proto.Unmarshal([]byte(pbSourceStr), circuitBreaker); err != nil {
		t.Fatal(err)
	}

	t.Log(circuitBreaker)

}
