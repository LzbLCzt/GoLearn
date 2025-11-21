package basic

import (
	"fmt"
	"testing"
)

func TestBasic2(t *testing.T) {
	defer func() {
		fmt.Println("defer")
	}()

}
