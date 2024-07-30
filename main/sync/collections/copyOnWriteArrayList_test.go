package collections

import (
	"fmt"
	"sync"
	"testing"
)

/*
todo 实现一个简单的CopyOnWriteArrayList
*/

type CopyOnWriteList struct {
	mutex sync.RWMutex
	slice []interface{}
}

func (c *CopyOnWriteList) Append(item interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newSlice := make([]interface{}, len(c.slice)+1)
	copy(newSlice, c.slice)
	newSlice[len(c.slice)] = item
	c.slice = newSlice
}

func (c *CopyOnWriteList) Get(idx int) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if idx < 0 || idx >= len(c.slice) {
		return nil, false
	}
	return c.slice[idx], true
}

func Test_CopyOnWriteArrayList(t *testing.T) {
	list := CopyOnWriteList{}
	list.Append("Hello")
	list.Append("World")

	if value, ok := list.Get(1); ok {
		fmt.Println("Found value:", value)
	}

}
