package context

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

/* Go 的 context 包用于管理跨 goroutine 的 生命周期控制 和 元数据传递，核心功能如下：
	1. 通知机制：通过取消信号（如超时、手动取消）通知相关 goroutine 停止工作
	2. 元数据传递：在调用链中传递请求范围的键值对数据（如请求 ID、用户身份）
 */

// -----------------------------------------
// todo 通知机制（取消信号）
func TestPractice(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go worker(ctx)
	time.Sleep(3 * time.Second)
	fmt.Println("main： 发送取消信号")
	cancel()

	// 等待 worker 退出
	time.Sleep(1 * time.Second)
}

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker: 收到取消信号，停止工作")
			return
		default:
			fmt.Println("Worker: 工作中...")
			time.Sleep(1 * time.Second)
		}
	}
}
// -----------------------------------------

// -----------------------------------------
// todo 通知机制（取消信号）
func TestPractice2(t *testing.T) {
	var wg sync.WaitGroup

	// todo 超时5s自动触发取消
	parentCtx, cancelParent := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelParent() // 确保资源释放

	// 启动 3 个任务，共享父 Context 的超时
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go longTask(parentCtx, i, &wg)
	}

	// todo 启动第 4 个任务，单独设置手动取消
	manualCtx, cancelManual := context.WithCancel(context.Background())
	wg.Add(1)
	go longTask(manualCtx, 4, &wg)

	// 模拟运行 3 秒后，手动取消第 4 个任务
	time.Sleep(3 * time.Second)
	fmt.Println("\n主程序: 手动取消任务 4")
	cancelManual()

	// 等待所有任务结束
	wg.Wait()
	fmt.Println("所有任务已终止")
}

func longTask(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取任务的截止时间（如果有）
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf("任务 %d 的截止时间: %v\n", id, deadline.Format("15:04:05"))
	}

	// 模拟工作循环
	for {
		select {
		case <-ctx.Done():
			// 检查取消原因
			err := ctx.Err()
			switch err {
			case context.Canceled:
				fmt.Printf("任务 %d 被手动取消\n", id)
			case context.DeadlineExceeded:
				fmt.Printf("任务 %d 因超时终止\n", id)
			}
			return
		default:
			fmt.Printf("任务 %d 工作中...\n", id)
			time.Sleep(1 * time.Second)
		}
	}
}
// -----------------------------------------

// -----------------------------------------
//todo 元数据传递（键值对）
type key string

const requestIDKey key = "request_id"

// 中间件：注入请求 ID 到 Context
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 生成请求 ID
		requestID := "123-abcd"
		// 将 requestID 存入 Context
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 处理函数：从 Context 读取请求 ID
func handler(w http.ResponseWriter, r *http.Request) {
	// 从 Context 提取 requestID
	ctx := r.Context()
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Fprintf(w, "Request ID: %s", requestID)
	} else {
		fmt.Fprint(w, "Request ID 不存在")
	}
}

func TestPractice3(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", middleware(http.HandlerFunc(handler)))
	http.ListenAndServe(":8080", mux)
}
// -----------------------------------------
//todo cancelCtx能保证父节点取消后，所有子节点自动触发取消（对应方法：propagateCancel）
/*
实现原理，通过ctx.Done（是一个阻塞的chan）确认节点是否cancel了
代码：
cancelCtx.propagateCancel()
	child.cancel() => cancelCtx.cancel()
		c.done.Store(closedchan) or close(d) (一旦chan关闭，则不再阻塞，不阻塞意味着收到了关闭信号)
 */
func TestPractice4(t *testing.T) {
	var wg sync.WaitGroup
	// 创建父 Context
	parentCtx, cancelParent := context.WithCancel(context.Background())

	// 创建两个子 Context
	childCtx1, _ := context.WithCancel(parentCtx)
	childCtx2, _ := context.WithCancel(parentCtx)

	// 启动监控协程
	wg.Add(2)
	go monitorChild(childCtx1, "子节点1", &wg)
	go monitorChild(childCtx2, "子节点2", &wg)

	// 模拟运行 2 秒后取消父 Context
	time.Sleep(2 * time.Second)
	fmt.Println("\n主程序: 取消父 Context")
	cancelParent()  // 触发父节点取消

	// 等待所有子节点退出
	wg.Wait()
	fmt.Println("所有子节点已终止")
}

// 监控子 Context 是否被取消
func monitorChild(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// 获取取消原因
			err := ctx.Err()
			if err == context.Canceled {
				fmt.Printf("[%s] 收到取消信号（原因: %v）\n", name, err)
			} else {
				fmt.Printf("[%s] 异常终止: %v\n", name, err)
			}
			return
		default:
			fmt.Printf("[%s] 运行中...\n", name)
			time.Sleep(1 * time.Second)
		}
	}
}
// -----------------------------------------
