package error

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

var DuplicateError = fmt.Errorf("duplicate error")
func TestError(t *testing.T) {
	err := returnErrorFunc()
	t.Logf("is duplicate error: %v", errors.Is(err, DuplicateError))
}

func returnErrorFunc() error {
	return DuplicateError
}