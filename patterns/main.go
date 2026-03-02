package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	// fetchParallel()
	paginate()
	fmt.Println("Time:", time.Since(start))
}

func paginate() {
	var finalResult []Transaction

	resultChannel := make(chan []Transaction)
	total := 100
	limit := 10

	for offset := 0; offset < total; offset += limit{
		go func(offset int) {
			resultChannel <- fetchPage(offset, limit)
		}(offset)
	}
	

	// query channel to get arrays and flatten it to one final array
	iterations := total/limit 
	// this part is important to block till our desired data is not retrieved
	for range iterations {
		val := <-resultChannel
		finalResult = append(finalResult, val...)
	}

	// if we want to implement timeout this can also be used:

	// select {
	// case val := <-resultChannel:
	// 	fmt.Println(val)
	// case <-time.After(100 * time.Millisecond):
	// 	fmt.Println("timeout")
	// 	break
		
	// }
	
	
	for _, val := range finalResult {
		fmt.Println(val.ID)
	}
	close(resultChannel)

}

func fetchParallel() {
	// user := fetchUser()
	// orders := fetchOrders()
	// payments := fetchPayments()

	userChannel := make(chan string)	
	orderChannel := make(chan string)
	paymentChannel := make(chan string)

	var wg sync.WaitGroup
	wg.Add(3)

	go func () {
		defer wg.Done()
		userChannel <- fetchUser()
	}()
	go func () {
		defer wg.Done()
		paymentChannel <- fetchPayments()
	}()
	go func () {
		defer wg.Done()
		orders := fetchOrders()
		orderChannel <- orders
	}()

	

	fmt.Println(<-userChannel, <-orderChannel, <-paymentChannel)
	wg.Wait()
}