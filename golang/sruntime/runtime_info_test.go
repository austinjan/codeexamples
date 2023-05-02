package sruntime

import (
	"runtime"
	"testing"
)

// numGoRoutines returns the number of goroutines
func numGoRoutines() int {
	return runtime.NumGoroutine()
}

func stackTrace() []byte {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	return buf
}

func MemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func TestRuntime(t *testing.T) {
	t.Log("Number of goroutines:", numGoRoutines())
	t.Log("Stack trace:", string(stackTrace()))
	t.Log("Allocated memory:", MemoryUsage())
}
