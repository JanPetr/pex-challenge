package metric

import (
	"sync"
)

// Dummy global metrics implementation
// Possible improvements - implement as a interface / struct with its variables, push metrics to any monitoring tool

var mu sync.Mutex
var metrics = make(map[string]int64)

func SetInt64(name string, v int64) {
	mu.Lock()
	defer mu.Unlock()

	metrics[name] = v
}

func AddInt64(name string, v int64) {
	mu.Lock()
	defer mu.Unlock()

	metrics[name] += v
}

func GetInt64(name string) int64 {
	return metrics[name]
}
