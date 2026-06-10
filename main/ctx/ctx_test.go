package ctx

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestContext(t *testing.T) {
	// 1. 创建名为 c1 的 context
	c1 := context.Background()

	// 2. 给 c1 添加一个 value
	type key string
	const keyUserID key = "user_id"
	c1 = context.WithValue(c1, keyUserID, "12345")

	// 3. 在 c1 基础上创建子 context
	c2 := context.WithValue(c1, key("name"), "Alice")
	c3 := context.WithValue(c1, "aaa", "bbb")

	// 验证 value 是否正确
	userID := c2.Value(keyUserID)
	t.Logf("user_id from c2: %v", userID)

	name := c2.Value(key("name"))
	t.Logf("name from c2: %v", name)

	aaa := c3.Value("aaa")
	t.Logf("aaa from c3: %v", aaa)
}

func TestContext2(t *testing.T) {
	c := context.Background()
	c = context.WithValue(c, "name", "lzb")
	_, ctx := errgroup.WithContext(c)

	v := ctx.Value("name")
	t.Logf("name: %v", v)
}

func TestContext3(t *testing.T) {
	c := context.Background()
	wait, egCtx := errgroup.WithContext(c)	// egCtx.Done 什么时候会有返回值（不再阻塞）：1 任意goroutine返回err 2 wait.Wait() 返回（所有goroutine执行完成了）
	egCtx = context.WithValue(egCtx, 0, "label 0")
	egCtx = context.WithValue(egCtx, 1, "label 1")
	egCtx = context.WithValue(egCtx, 2, "label 2")

	for i := 0; i < 3; i++ {
		tmp := i
		wait.Go(func() error {
				if tmp == 1 {
					return errors.New("error")
				}
				for {
					t.Logf("i: %v, ctx value: %v", tmp, egCtx.Value(tmp))
					time.Sleep(time.Second * 5)
				}
		})
	}

	t.Logf("start waiting")
	if err := wait.Wait(); err != nil {
		t.Logf("wait error: %v", err)
	}

	t.Logf("wait done")
}


