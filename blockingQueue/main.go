package main

import (
	"fmt"
	"sync"
	"time"
)


func main() {
	fmt.Println("Hello vi")

	bq := NewBlockingQueue(1)
	var wg sync.WaitGroup

	wg.Add(2)

	// Producer
	go func() {
		defer wg.Done()
		for i := 1; i < 17; i++ {
			bq.Put(i)
			time.Sleep(200 * time.Millisecond)
		}
		
	}()

	// Consumer
	go func() {
		defer wg.Done()
		// time.Sleep(1 * time.Second)
		for range 4 {
			for range 4 {
				val, _ := bq.Remove()
				fmt.Println("Consumed:", val)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()
}