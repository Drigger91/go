package main

import (
	"fmt"
	"sync"
)


type BlockingQueue struct {
	queue []int
	cond *sync.Cond
	size int
	closed bool
}

func NewBlockingQueue(size int) *BlockingQueue {
	return &BlockingQueue{
		queue: make([]int, 0, size),
		cond: sync.NewCond(&sync.Mutex{}),
		size: size,
	}
} 

func (q *BlockingQueue) Print() {
	for _, val := range q.queue {
		fmt.Print(val, " ")
	}
	fmt.Println()
}

func (q *BlockingQueue) checkCloseSanity() error {
	if q.closed {
		return fmt.Errorf("Queue already closed!!")
	}
	return nil
}

func (q *BlockingQueue) Put(val int) error {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if err := q.checkCloseSanity(); err != nil {
		return err
	}	
	for len(q.queue) == q.size {
		// wait till someone removes
		if q.closed {
			return fmt.Errorf("Queue closed")
		}
		fmt.Println("Queue size full, wating for removal...")
		q.cond.Wait()
	}
	
	q.queue = append(q.queue, val)
	
	q.cond.Signal()
	return nil
}

func (q *BlockingQueue) Remove() (int, error) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if err := q.checkCloseSanity(); err != nil {
		return 0, err
	}
	
	for len(q.queue) == 0 {
		if q.closed {
			return -1, fmt.Errorf("Queue closed")
		}
		q.cond.Wait()
		fmt.Println("No value in queue waiting...")
	}
	topVal := q.queue[0]
	q.queue = q.queue[1:]
	q.cond.Signal()
	return topVal, nil
}

func(q *BlockingQueue) Close() (bool, error) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if err := q.checkCloseSanity(); err != nil {
		return false, err
	}
	fmt.Println("Closing queue final state: ")
	// q.Print()
	q.closed = true
	q.cond.Broadcast()
	return true, nil
}
