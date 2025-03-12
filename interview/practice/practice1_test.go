package practice

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var mu1 sync.Mutex
var chain string

// todo deadlock occurs
func TestPractice1(t *testing.T) {
	chain = "main"
	A()
	fmt.Println(chain)
}
func A() {
	mu1.Lock()
	defer mu1.Unlock()
	chain = chain + " --> A"
	B()
}
func B() {
	chain = chain + " --> B"
	C()
}
func C() {
	mu1.Lock()
	defer mu1.Unlock()
	chain = chain + " --> C"
}

//-----------------------------------------------------------

// todo deadlock occurs why?
/*
todo 锁升级问题：当一个线程（或goroutine）持有读锁，并且另一个线程尝试获取写锁时，后者必须等待所有读锁释放。如果在这种情况下，持有读锁的线程再次尝试获取读锁（即使是可重入的读锁），它也会被阻塞，因为系统可能已经将写锁请求排在了前面，这就形成了一个循环等待的情况，即死锁。
todo 重入问题：虽然sync.RWMutex在Go中支持同一个goroutine对读锁的重入，但在写锁请求已经挂起的情况下，再次尝试获取读锁会导致死锁，因为写锁请求会阻止进一步的读锁授予，即使是给已经持有读锁的goroutine。
*/
var mu sync.RWMutex
var count int

func TestPractice2(t *testing.T) {
	go A1()
	time.Sleep(2 * time.Second)
	mu.Lock()
	defer mu.Unlock()
	count++
	fmt.Println(count)
}
func A1() {
	mu.RLock()
	defer mu.RUnlock()
	B2()
}
func B2() {
	time.Sleep(5 * time.Second)
	C2()
}
func C2() {
	mu.RLock()
	defer mu.RUnlock()
}

// todo panic: WaitGroup is reused before previous Wait has returned
func TestPractice3(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond)
		wg.Done()
		wg.Add(1)
	}()
	wg.Wait()
}

// -----------------------------------------------------------
// todo 实现一个类似sync.Once的单例
// todo 问题：f函数可能执行多次，因为goroutine1更新done = 1后，goroutine2缓存的可能是旧值（done=0），因此goroutine获取lock之后判断o.done==0通过，会再次执行f函数。
type Once struct {
	m    sync.Mutex
	done uint32
	//done atomic.Uint32
}

func (o *Once) Do(f func()) {
	if o.done == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		o.done = 1
		f()
	}
}

func TestPractice4(t *testing.T) {
	o := Once{m: sync.Mutex{}, done: 0}
	for i := 0; i < 10000; i++ {
		go o.Do(func() {
			fmt.Println("once")
		})
	}
	time.Sleep(5 * time.Second)
}

// todo 修复 once.Do() 的问题
type Once2 struct {
	m    sync.Mutex
	done uint32
}

func (o *Once2) Do(f func()) {
	// 使用 atomic.LoadUint32 保证并发安全地读取 done 的值
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()

	// 在获取锁后再次检查 done，防止多个 goroutine 进入临界区
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1) // 原子更新 done 标志
		f()
	}
}

func TestPractice4_repair(t *testing.T) {
	o := Once2{}
	for i := 0; i < 10000; i++ {
		go o.Do(func() {
			fmt.Println("once")
		})
	}
	time.Sleep(5 * time.Second)
}

// -----------------------------------------------------------
// todo all goroutines are asleep - deadlock
// todo sync.Mutex 必须始终通过指针来使用，因为它的内部状态不应被复制。当你复制一个包含 Mutex 的结构体时，
// todo Mutex 的内部状态也被复制，这会导致运行时错误（panic），因为锁的状态（是否被锁定等）在复制后变得不一致
// todo 结论：sync.Mutex的用法不对，不应该被复制后使用
type MyMutex struct {
	count int
	sync.Mutex
}

func TestPractice5(t *testing.T) {
	var mu MyMutex
	mu.Lock()
	var mu2 = mu
	mu.count++
	mu.Unlock()
	mu2.Lock()
	mu2.count++
	mu2.Unlock()
	fmt.Println(mu.count, mu2.count)
}

//-----------------------------------------------------------
