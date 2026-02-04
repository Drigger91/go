package l2

import (
	"fmt"
	"time"
)

func RateLimiter() {
	fmt.Println("Inside rate limiter")
	jobCh := make(chan int)

	tc := time.NewTicker(500 * time.Millisecond);
	go consumeJobsWithLimit(tc, jobCh)

	for i := range(20) {
		jobCh <- i
	}
	close(jobCh) // before waiting we have to close the channel
}

func consumeJobsWithLimit(tc *time.Ticker, jobCh <- chan int) {
	defer tc.Stop()
	// ideally this should be added as for select as new cases will be easier to add
	for {
		select {
		case <-tc.C:
			select {
			case job, ok := <-jobCh:
				if !ok {
					return
				}
				fmt.Println(job)
			default:
				fmt.Println("no job left to processed")
			}
		} 
	}
}

