package sync

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type ConnectionManager struct {
	addresses    []string
	addressesMux sync.RWMutex
}

type ConnectionManager2 struct {
	addresses atomic.Value
}

func TestSync(t *testing.T) {
	now := time.Now()
	time.Sleep(time.Second)
	duration := time.Since(now)

	fmt.Println(int64(duration))
}
