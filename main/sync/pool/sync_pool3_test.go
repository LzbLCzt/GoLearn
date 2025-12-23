package sync

import (
	"sync"
	"testing"
)

/*
	复用结构题
*/

type TempStruct struct {
	Data  []byte
	Count int
}

var structPool = sync.Pool{
	New: func() interface{} { return &TempStruct{} },
}

func TestSyncPool3(t *testing.T) {
	obj := structPool.Get().(*TempStruct)
	defer structPool.Put(obj)

	// 使用前重置对象状态
	obj.Data = obj.Data[:0]
	obj.Count = 0

	// 使用对象...
	obj.Data = append(obj.Data, "new data"...)
	obj.Count++
}
