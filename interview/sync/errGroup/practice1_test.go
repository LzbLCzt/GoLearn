package practice

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
	"testing"
	"time"
)

//-----------------------------------------------------------
//todo errGroup可以收集一组goroutine中某个goroutine的错误
func TestPractice1(t *testing.T) {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.dfasdfadfads.com/",
	}
	g, ctx := errgroup.WithContext(context.Background())

	var result sync.Map

	for _, url := range urls {
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return err
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			result.Store(url, resp.StatusCode)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("error: ", err)
	}

	result.Range(func(key, value interface{}) bool {
		fmt.Printf("fetch url: %s, status: %d\n", key, value)
		return true
	})
}
//-----------------------------------------------------------
//todo 利用errGroup的WithContext，某个goroutine失败时，可以取消所有未完成的其他goroutine
//todo errGroup的WithContext能够支持取消其他未完成的goroutine，但是取消操作需要自己通过代码实现
func TestPractice2(t *testing.T) {
	g, ctx := errgroup.WithContext(context.Background())

	for i := 1; i <= 3; i++ {
		i := i
		g.Go(func() error {
			return worker(ctx, i)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("all workers finished successfully")
	}
}

func worker(ctx context.Context, id int) error {
	fmt.Printf("Worker %d started\n", id)
	select {
	case <-time.After(time.Duration(id) * time.Second): // 模拟不同工作时间
		fmt.Printf("Worker %d finished\n", id)
	case <-ctx.Done(): // 检查是否被取消
		fmt.Printf("Worker %d cancelled\n", id)
		return ctx.Err()
	}
	if id == 2 {
		return fmt.Errorf("worker %d failed", id) // 模拟错误
	}
	return nil
}
//-----------------------------------------------------------
