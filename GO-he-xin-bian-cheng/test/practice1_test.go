package test

import "testing"

func TestPractice1(t *testing.T) {
	c := make(chan int, 1)
	go func(c chan int) {
		for {
			select {
			case c <- 0:
			case c <- 1:
			}
		}
	}(c)

	for i := 0; i < 10; i++ {
		t.Log(<-c)
	}

	close(c)
	t.Log("aa")
	t.Log(<-c)
	t.Log(<-c)
	t.Log(<-c)
	t.Log(<-c)
	t.Log(<-c)
}
