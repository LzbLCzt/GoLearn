package practice

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//var mu sync.Mutex
//var chain string

// todo deadlock occurs
//func TestPractice1(t *testing.T) {
//	chain = "main"
//	A()
//	fmt.Println(chain)
//}
//func A() {
//	mu.Lock()
//	defer mu.Unlock()
//	chain = chain + " --> A"
//	B()
//}
//func B() {
//	chain = chain + " --> B"
//	C()
//}
//func C() {
//	mu.Lock()
//	defer mu.Unlock()
//	chain = chain + " --> C"
//}

//-----------------------------------------------------------

// todo deadlock occurs why?
/*
todo 锁升级问题：当一个线程（或goroutine）持有读锁，并且另一个线程尝试获取写锁时，后者必须等待所有读锁释放。如果在这种情况下，持有读锁的线程再次尝试获取读锁（即使是可重入的读锁），它也会被阻塞，因为系统可能已经将写锁请求排在了前面，这就形成了一个循环等待的情况，即死锁。
todo 重入问题：虽然sync.RWMutex在Go中支持同一个goroutine对读锁的重入，但在写锁请求已经挂起的情况下，再次尝试获取读锁会导致死锁，因为写锁请求会阻止进一步的读锁授予，即使是给已经持有读锁的goroutine。
*/
var mu sync.RWMutex
var count int

func TestPractice1(t *testing.T) {
	go A()
	time.Sleep(2 * time.Second)
	mu.Lock()
	defer mu.Unlock()
	count++
	fmt.Println(count)
}
func A() {
	mu.RLock()
	defer mu.RUnlock()
	B()
}
func B() {
	time.Sleep(5 * time.Second)
	C()
}
func C() {
	mu.RLock()
	defer mu.RUnlock()
}
