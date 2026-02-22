package numberprinter

import (
	"fmt"
	"sync"
)

// print odd from one thread even from another thread
func Print() {
	n := 12
	var wg sync.WaitGroup
	wg.Add(2)
	// loop ->	starts from 0 -> push 0 -> push 1 ....
	// odd thread, even thread
	// when odd thread will push it'll wait for nextEven to be printed
	// evenPush -> waitForOddRelease -> oddPush -> do OddRelease -> wait for evenRelease -> evenPush

	turn := 0 
	cond := sync.NewCond(&sync.RWMutex{})

	// Even producer
	go func() {
		defer wg.Done()
		for i := 0; i <= n; i += 2 {
			cond.L.Lock()
			// Always check the state in a loop before waiting
			for turn != 0 {
				cond.Wait()
			}
			fmt.Println(i)
			turn = 1         
			cond.Signal()   
			cond.L.Unlock()
		}
	}()

	// Odd producer
	go func() {
		defer wg.Done()
		for i := 1; i <= n; i += 2 {
			cond.L.Lock()
			for turn != 1 {
				cond.Wait()
			}
			fmt.Println(i)
			turn = 0        
			cond.Signal()  
			cond.L.Unlock()
		}
	}()
	

	wg.Wait()
}


// channel approach is better than cond approach as we should not communicate by sharing memory we should share
// memory by communicating. (in prev approach we were relying on turn variable)
func PrintChannel() {
	var wg sync.WaitGroup
	wg.Add(2)
	n := 12
	evenCh := make(chan bool)
	oddCh := make(chan bool)
	go func ()  {
		defer wg.Done()
		for i := 0; i <= n; i+=2 {
			// block till allowed
			<-evenCh
			fmt.Println(i)

			if (i < n) {
				oddCh <- true
			}
		}
	}()

	go func ()  {
		defer wg.Done()
		for i := 1; i <= n; i+=2 {
			// block till allowed
			<-oddCh
			fmt.Println(i)

			if (i < n) {
				evenCh <- true
			}
		}
	}()

	evenCh <- true
	wg.Wait()
}