package main

import (
	"fmt"
	"time"
)

func main() {
	// @doc https://stackoverflow.com/a/35009735
	// ticker1 := time.NewTicker(24 * time.Hour)  // This should be under production
	ticker1 := time.NewTicker(10 * time.Second)

	for c := range ticker1.C {
		fmt.Println(c.Format("2006/01/02 15:00:00"))
	}
}
