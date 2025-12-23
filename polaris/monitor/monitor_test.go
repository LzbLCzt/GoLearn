package monitor

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMonitor(t *testing.T) {
	cpu := runtime.NumCPU()
	procs := runtime.GOMAXPROCS(0)
	fmt.Println("cpu num:", cpu, " GOMAXPROCS:", procs)
}
