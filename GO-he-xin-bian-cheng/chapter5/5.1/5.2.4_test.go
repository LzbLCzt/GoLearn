package __1

import (
	"fmt"
	"sync"
	"testing"
)

func Test5_2_4(t *testing.T) {
	taskCh := make(chan task, 10)
	resultCh := make(chan int, 10)

	wg := &sync.WaitGroup{}
	go InitTask(taskCh, resultCh, 100)

	go DistributeTask(taskCh, wg, resultCh)

	res := ProcessResult(resultCh)
	fmt.Println(res)
}

func DistributeTask(taskchan <-chan task, wait *sync.WaitGroup, result chan int) {
	for task := range taskchan {
		wait.Add(1)
		go ProcessTask(task, wait)
	}
	wait.Wait()
	close(result)
}

func ProcessTask(t task, result *sync.WaitGroup) {
	t.do()
	result.Done()
}

func InitTask(taskchan chan<- task, r chan int, p int) {
	qu := p / 10
	mod := p % 10
	high := qu * 10
	for j := 0; j < qu; j++ {
		b := 10*j + 1
		e := 10 * (j + 1)
		tsk := task{
			begin:  b,
			end:    e,
			result: r,
		}
		taskchan <- tsk
	}
	if mod != 0 {
		tsk := task{
			begin:  high + 1,
			end:    p,
			result: r,
		}
		taskchan <- tsk
	}

	close(taskchan)
}

type task struct {
	begin, end int
	result     chan<- int
}

func (t *task) do() {
	sum := 0
	for i := t.begin; i <= t.end; i++ {
		sum += i
	}
	t.result <- sum
}

// 读取结果通道，汇总结果
func ProcessResult(resultchan chan int) int {
	sum := 0
	for r := range resultchan {
		sum += r
	}
	return sum
}
