package goroutine

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

/*
goroutine设置超时机制
*/
func TestGoRoutineOvertime2(t *testing.T) {
	total := 12
	var num int32
	log.Println("begin")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	for i := 0; i < total; i++ {
		go func() {
			//time.Sleep(3 * time.Second)
			atomic.AddInt32(&num, 1)
			if atomic.LoadInt32(&num) == 10 {
				cancelFunc()
			}
		}()
	}
	for i := 0; i < 5; i++ {
		go func() {

			select {
			case <-ctx.Done():
				log.Println("ctx1 done", ctx.Err())
			}

			for i := 0; i < 2; i++ {
				go func() {
					select {
					case <-ctx.Done():
						log.Println("ctx2 done", ctx.Err())
					}
				}()
			}

		}()
	}

	time.Sleep(time.Second*5)
	log.Println("end", ctx.Err())
	fmt.Printf("执行完毕 %v", num)
}
