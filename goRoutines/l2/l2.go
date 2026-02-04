package l2

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func L2Main() {
	fmt.Println("Hello from l2")
	//consumeJobs()
	oddCh := make(chan int)
	evenCh := make(chan int)
	mergedCh := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go producerOdd(oddCh, &wg)
	go producerEven(evenCh, &wg)

	// merger
	go func() {
		defer close(mergedCh)

		oddOpen, evenOpen := true, true

		for oddOpen || evenOpen {
			select {
			case v, ok := <-oddCh:
				if !ok {
					oddOpen = false
					continue
				}
				mergedCh <- v

			case v, ok := <-evenCh:
				if !ok {
					evenOpen = false
					continue
				}
				mergedCh <- v
			}
		}
	}()

	for v := range mergedCh {
		fmt.Println("val", v)
	}

	wg.Wait()
}

func producerOdd(ch chan int, wg *sync.WaitGroup) {
	for i := 0; i <= 12; i+=2 {
		ch <- i
	}
	close(ch)
	wg.Done()
}
func producerEven(ch chan int, wg *sync.WaitGroup) {
	for i := 1; i <= 11; i+=2 {
		ch <- i
	}
	close(ch)
	wg.Done()
}
func consumeJobs() {
	jobCh := make(chan int)
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 1; i<=3; i++ {
		go func() {
			for job := range(jobCh) {
				workerFunc(strconv.Itoa(i), job)
			}
			wg.Done()
		}()
	}
	
	for i := range(20) {
		jobCh <- i
		
	}
	close(jobCh) // before waiting we have to close the channel
	wg.Wait()
}

func consumeJobsWithCancellation() {
	jobCh := make(chan int)
	doneCh := make(chan struct{})

	var wg sync.WaitGroup
	workerCount := 3
	wg.Add(workerCount)

	for i := 1; i <= workerCount; i++ {
		id := i

		go func() {
			defer wg.Done()

			for {
				select {
				case job, ok := <-jobCh:
					if !ok {
						return
					}
					workerFunc(strconv.Itoa(id), job)

				case <-doneCh:
					fmt.Printf("Worker %d cancelled\n", id)
					return
				}
			}
		}()
	}
	go func() {
		for i := 0; i < 20; i++ {
			jobCh <- i
		}
		close(jobCh)
	}()

	// cancel after 1 second
	time.Sleep(1 * time.Second)
	close(doneCh)

	wg.Wait()
}
func workerFunc(id string, job int) {
	time.Sleep(time.Duration(1 * time.Second))
	fmt.Printf("Worker %s processed job %d\n", id, job)
}