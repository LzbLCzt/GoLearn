package practice

import (
	"fmt"
	"testing"
)

func SendNotification(user string) chan string {
	ch := make(chan string, 500)
	go func() {
		ch <- fmt.Sprintf("Hi %s, welcome to our site!", user)
	}()
	return ch
}

func Test_notification1(t *testing.T) {
	barry := SendNotification("barry")
	shirdon := SendNotification("shirdon")

	fmt.Println(<-barry)
	fmt.Println(<-shirdon)
}
