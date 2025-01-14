package __1

import (
	"context"
	"fmt"
	"testing"
	"time"
)

//todo context

func Test_5_3_5_A(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go work(ctx, "work1")

	ctxWithDeadline, _ := context.WithDeadline(ctx, time.Now().Add(3*time.Second)) //todo 3s后自动退出
	go work(ctxWithDeadline, "work2")

	oc := otherContext{ctx}
	valueCtx := context.WithValue(oc, "key", "god andes,pass from main")
	go workWithValue(valueCtx, "work3")
	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(5 * time.Second)
	fmt.Println("main stop")

}

type otherContext struct {
	context.Context
}

func work(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s get msg to cancel\n", name)
			return
		default:
			fmt.Printf("%s is running \n", name)
			time.Sleep(1 * time.Second)
		}
	}
}

func workWithValue(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s get msg to cancel\n", name)
			return
		default:
			value := ctx.Value("key").(string)
			fmt.Printf("%s is running value=%s \n", name, value)
			time.Sleep(1 * time.Second)
		}
	}
}
