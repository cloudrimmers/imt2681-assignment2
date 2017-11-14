package benchmark

import (
	"log"
	"time"

	"github.com/cloudrimmers/imt2681-assignment3/lib/reflectUtil"
)

//Benchmark - Duration in ms a function takes, to use me, defer me with time.Now()
func Benchmark(start time.Time) {
	log.Printf("Duration of %v(...) : %dms", reflectUtil.GetCallerName(), time.Now().Nanosecond()/1e6-start.Nanosecond()/1e6)
}
