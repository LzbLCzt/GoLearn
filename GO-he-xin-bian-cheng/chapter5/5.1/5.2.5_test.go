package __1

import (
	"fmt"
	"testing"
	"time"
)

//todo future模式

func Test_5_2_5_A(t *testing.T) {
	q := query{sql: make(chan string, 1), result: make(chan string, 1)}
	go execQuery(q)

	q.sql <- "select * from user"
	time.Sleep(time.Second)
	fmt.Println(<-q.result)
}

type query struct {
	sql    chan string
	result chan string
}

func execQuery(q query) {
	go func() {
		sql := <-q.sql
		q.result <- "result from " + sql
	}()
}