package benchmark

import (
	"testing"
	"time"
)

func TestBenchmark(t *testing.T) {
	defer Benchmark(time.Now())
}
