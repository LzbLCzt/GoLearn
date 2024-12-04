package goroutine

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

/*
todo Go语言高性能协程池实现与原理解析
*/

func TestPool(t *testing.T) {
	pool := NewPool(10)
	pool.Start()
	defer pool.Stop()
	for i := 0; i < 5; i++ {
		err := pool.Submit(func() error {
			fmt.Printf("开始处理第 %d 个task\n", i)
			time.Sleep(1 * time.Millisecond)
			fmt.Printf("处理完成第 %d 个task\n", i)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	//time.Sleep(10 * time.Second)
}

type Task struct {
	Handler func() error
	Result  chan error
}

type Pool struct {
	capacity    int
	active      int
	tasks       chan *Task //任务队列
	quit        chan bool
	workerQueue chan *worker //工作协程队列
	mutex       sync.Mutex
}

type worker struct {
	pool *Pool
}

func NewPool(capacity int) *Pool {
	if capacity <= 0 {
		capacity = 1
	}

	return &Pool{
		capacity:    capacity,
		tasks:       make(chan *Task, capacity*2),
		quit:        make(chan bool),
		workerQueue: make(chan *worker, capacity),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.capacity; i++ {
		w := &worker{pool: p}
		p.workerQueue <- w
		p.active++
		go w.run() //启动数量 = capacity的goroutine,等待处理任务
	}
}

func (w *worker) run() {
	for {
		select {
		case task := <-w.pool.tasks:
			if err := task.Handler(); err != nil {
				task.Result <- err
			} else {
				task.Result <- nil
			}
			//工作完成后，将自己放回队列
			w.pool.workerQueue <- w
		case <-w.pool.quit:
			return
		}
	}
}

// 创建一个task并执行
func (p *Pool) Submit(handler func() error) error {
	task := &Task{
		Handler: handler,
		Result:  make(chan error, 1),
	}

	p.tasks <- task
	return <-task.Result
}

func (p *Pool) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.active > 0 {
		close(p.quit)
		p.active = 0
	}
}

// todo 为了提升协程池的性能，来个任务批处理功能
// -------------------------------------------------------
type BatchPool struct {
	*Pool
	batchSize int
	batchChan chan []*Task
}

func (bp *BatchPool) processBatch() {
	batch := make([]*Task, 0, bp.batchSize) //切片的初始长度 = 0,切片的初始容量 = bp.batchSize(容量指的是在需要重新分配内存之前，切片可以存储的元素数量的上限)
	timer := time.NewTimer(100 * time.Millisecond)
	for {
		select {
		case task := <-bp.tasks:
			batch = append(batch, task)
			if len(batch) >= bp.batchSize {
				bp.batchChan <- batch
				batch = make([]*Task, 0, bp.batchSize)
			}

		case <-timer.C:
			if len(batch) > 0 {
				bp.batchChan <- batch
				batch = make([]*Task, 0, bp.batchSize)
			}
			timer.Reset(100 * time.Millisecond)
		}
	}
}

//-------------------------------------------------------

// todo 性能测试
// -------------------------------------------------------
func BenchmarkPool(b *testing.B) {
	pool := NewPool(runtime.NumCPU())
	pool.Start()
	defer pool.Stop()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.Submit(func() error {
				time.Sleep(time.Second)
				return nil
			})
		}
	})
}

//-------------------------------------------------------
