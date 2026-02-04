package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Drigger91/go/concurrency/l2"
)

func main() {
	// ch := make(chan string)
	// go func(){
	// 	time.Sleep(1 * time.Second)
	// 	ch <- "Hello from routine\n"
	// }()
	// fmt.Println(<-ch)
	// fmt.Println("Hello from main")

	// intCh := make(chan int, 2)

	// go func(){
	// 	for i := 0; i < 20; i++ {
	// 		fmt.Println("Sending", i)
	// 		intCh <- i
	// 		fmt.Println("Sent", i)
	// 	}
	// 	close(intCh)
	// }()

	// for val := range intCh {
	// 	fmt.Println(val)
	// }

	// // closing buffered channel is imp
	// fmt.Println("Int channel")
	// intCh = make(chan int)
	// go func(){
	// 	for i := 0; i < 5; i++ {
	// 		intCh <- i
	// 	}
	// 	close(intCh)
	// }()
	
	// for i := range(intCh) {
	// 	fmt.Println(i)
	// }
	// Level1
	// ch := make(chan bool)

	// go func(){
	// 	<-ch
	// 	fmt.Println("Start")
	// }()

	// time.Sleep(500 * time.Microsecond)
	// ch <- true
	// fmt.Println("Done")

	// channelCancelMain()
	// channelCancelWg()

	// waitGroupDemo()

	// ch := make(chan string)
	// go func() {
	// 	defer close(ch)
    //     msgs := []string{"msg1", "msg2", "msg3"}
    //     for _, m := range msgs {
    //         time.Sleep(400 * time.Millisecond)
    //         ch <- m
    //     }
    // }()

    // waitForMessage(ch, 1 * time.Second)
    // fmt.Println("Main done")
	//l2.L2Main()
	l2.RateLimiter()
}

func channelCancelMain() {
	ch := make(chan bool)
	go func(){
		for{
			select{
			case <-ch:
				fmt.Println("Stopped")
				return
			default:
				fmt.Println("start..")
				time.Sleep(300 * time.Millisecond)
			}
		}
		
	}()
	time.Sleep(1 * time.Second)
	close(ch)
	time.Sleep(100 * time.Millisecond) // this is needed or wgGroup
	fmt.Println("done")
}

func channelCancelWg(){
	ch := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1) // buffer is neccessary
	go func(){
		defer wg.Done()
		for{
			select{
			case <-ch:
				fmt.Println("Stopped")
				return
			default:
				fmt.Println("start..")
				time.Sleep(300 * time.Millisecond)
			}
		}
		
	}()
	time.Sleep(1 * time.Second)
	close(ch)
	wg.Wait()
	fmt.Println("done")
}

func waitGroupDemo() {
	var wg sync.WaitGroup

	wg.Add(2) // arg: delta -> this means it will wait till 2 wg.Done()
	// If delta > all go routines then deadlock happens (i.e all goroutines are sleep error handle with caution)
	go func(){
		fmt.Println("First")
		wg.Done()
	}()
	go func(){
		fmt.Println("Second")
		wg.Done()
	}()
	go func(){
		fmt.Println("Third")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Completed")
}
func waitForMessage(ch <-chan string, timeout time.Duration) {
    deadline := time.After(timeout)
    for {
        select {
        case msg, ok := <-ch:
            if !ok {
                fmt.Println("Channel closed. Stopping...")
				return
            }
            fmt.Println("Received:", msg)

        case <-deadline:
            fmt.Println("Total timeout reached. Exiting...")
        }
    }
}
