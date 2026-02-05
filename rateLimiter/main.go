package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Rate limiter")
	// rl := NewBasicRateLimiter(5, 1*time.Second)

	// for range 9 {
	// 	fmt.Println(rl.Allow(1))
	// }
	// time.Sleep(2 * time.Second)
	// fmt.Println("-----------")
	// for range 5 {
	// 	fmt.Println(rl.Allow(1))
	// }

	rl := NewBucketRateLimiter(5, 1*time.Second)

	for range 9 {
		fmt.Println(rl.Allow(1))
	}
	time.Sleep(2 * time.Second)
	fmt.Println("-----------")
	for range 5 {
		fmt.Println(rl.Allow(1))
	}
}
