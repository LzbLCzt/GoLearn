package sync

import (
	"fmt"
	"sync"
	"testing"
)

var wg sync.WaitGroup
var once sync.Once

func func1(ch1 chan<- int) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		fmt.Println("ch1接受元素:", i)
		ch1 <- i
	}
	close(ch1)
}

func func2(ch1 <-chan int, ch2 chan<- int, tag int) {
	defer wg.Done()
	for {
		x, ok := <-ch1
		if !ok {
			break
		}
		fmt.Println("ch1通道发送元素给ch2：", x, "tag: ", tag)
		ch2 <- 2 * x
	}

	once.Do(func() {
		fmt.Println("ch2通道关闭, ", "tag: ", tag)
		close(ch2)
	})
}

func Test_2(t *testing.T) {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)

	wg.Add(3)

	go func1(ch1)
	go func2(ch1, ch2, 1)
	go func2(ch1, ch2, 2)

	wg.Wait()

	for ret := range ch2 {
		fmt.Println("ch2通道发送元素: ", ret)
	}
}

func Test_3(t *testing.T) {
	ch := make(chan int, 10)
	go func() {
		/*todo 什么情况下ok返回的是false:通道 ch 已经被关闭，并且通道中没有剩余的数据可供接收时；
		todo 具体来说，当通道关闭时，可以继续从通道接收数据，直到通道中的所有已发送的数据都被接收完毕。一旦通道为空，任何进一步的接收操作将不会得到数据，此时接收操作返回的 ok 将是 false
		*/
		x, ok := <-ch
		if !ok {
			fmt.Println("ch通道关闭")
		}
		fmt.Println(x)
	}()
}

func Test_4(t *testing.T) {
	once := sync.Once{}
	for i := 0; i < 1000; i++ {
		once.Do(func() {
			fmt.Println("do only once")
		})
	}
}
