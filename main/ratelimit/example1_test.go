package ratelimit

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

// todo time.Ticker 可以用来实现基于时间间隔的简单限流。这种方法可以控制事件的发生频率，例如，限制处理请求的速率
func TestExample1(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 10; i++ {
		<-ticker.C
		fmt.Println("Tick at", time.Now())
	}
}

// todo golang.org/x/time/rate 包提供了一个更为复杂和实用的限流器，支持令牌桶算法
// todo 在这个例子中，rate.NewLimiter(1, 3) 创建了一个新的限流器，每秒生成一个令牌，桶的容量为3。这意味着即使在开始时可以快速处理三个请求，但之后的请求将被限制为每秒一个
func TestExample2(t *testing.T) {
	limiter := rate.NewLimiter(1, 3) // 每秒生成1个令牌，桶大小为3
	fmt.Println("limiter:", limiter)
	for i := 0; i < 10; i++ {
		ctx := context.Background()
		if err := limiter.Wait(ctx); err != nil {
			fmt.Println("Rate limit error:", err)
			return
		}
		fmt.Println("Request processed at", time.Now())
	}
}
