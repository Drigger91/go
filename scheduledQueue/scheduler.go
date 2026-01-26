package main

import (
	"container/heap"
	"fmt"
	"log"
	"sync"
	"time"
)

// Scheduler manages scheduled tasks and executes them at their scheduled time
type Scheduler struct {
	mu          sync.RWMutex
	taskHeap    TaskHeap
	taskMap     map[string]*Task // taskId -> task for O(1) lookup
	stopCh      chan struct{}
	doneCh      chan struct{}
	started     bool
	stopped     bool
	wg          sync.WaitGroup
	nextTaskCh  chan struct{} // signal when new task is added
}

// NewScheduler creates a new scheduler instance
func NewScheduler() *Scheduler {
	return &Scheduler{
		taskHeap:   make(TaskHeap, 0),
		taskMap:    make(map[string]*Task),
		stopCh:     make(chan struct{}),
		doneCh:     make(chan struct{}),
		nextTaskCh: make(chan struct{}, 1),
	}
}

// Start starts the scheduler background loop
// Calling Start() multiple times should not create multiple loops
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		fmt.Println("Scheduler already started")
		return
	}

	s.started = true
	s.stopped = false

	heap.Init(&s.taskHeap)

	s.wg.Add(1)
	go s.run()
}

// run is the main scheduler loop that executes tasks
func (s *Scheduler) run() {
	defer s.wg.Done()
	defer close(s.doneCh)

	var timer *time.Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		s.mu.Lock()
		
		// Wait for tasks if heap is empty
		for s.taskHeap.Len() == 0 {
			s.mu.Unlock()
			select {
			case <-s.stopCh:
				return
			case <-s.nextTaskCh:
				s.mu.Lock()
				continue
			}
		}

		// Get the next task (peek at heap top)
		nextTask := s.taskHeap[0]
		now := time.Now()

		// Check if task was cancelled
		if _, exists := s.taskMap[nextTask.id]; !exists {
			heap.Pop(&s.taskHeap)
			s.mu.Unlock()
			continue
		}

		if now.Before(nextTask.runAt) {
			// Task is not ready yet, wait until it's time
			waitDuration := nextTask.runAt.Sub(now)
			s.mu.Unlock()

			// Stop previous timer if exists
			if timer != nil {
				timer.Stop()
			}
			timer = time.NewTimer(waitDuration)

			select {
			case <-s.stopCh:
				return
			case <-timer.C:
			case <-s.nextTaskCh:
				continue
			}
		} else {
			s.mu.Unlock()
		}

		// Pop and execute the task
		s.mu.Lock()
		if s.taskHeap.Len() == 0 {
			s.mu.Unlock()
			continue
		}

		// Pop the task from heap
		nextTask = heap.Pop(&s.taskHeap).(*Task)
		
		// Check if task was cancelled (not in map anymore)
		if _, exists := s.taskMap[nextTask.id]; !exists {
			// Task was cancelled, skip it
			s.mu.Unlock()
			continue
		}

		// Remove from map before unlocking
		delete(s.taskMap, nextTask.id)
		s.mu.Unlock()

		// Execute task handler (without lock)
		s.executeTask(nextTask)
	}
}

// executeTask executes a task and handles errors
func (s *Scheduler) executeTask(task *Task) {
	now := time.Now()
	fmt.Printf("Executed task=%s payload=%s at=%s\n", task.id, task.payload, now.Format(time.RFC3339))

	if task.handler != nil {
		// If handler returns an error, we log it but continue
		// For now, handler is a func(), but we could extend it
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Task %s panicked: %v", task.id, r)
				}
			}()
			task.handler()
		}()
	}
}

// Schedule schedules a task to run after the specified delay
func (s *Scheduler) Schedule(taskId string, delay time.Duration, payload string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return fmt.Errorf("scheduler is stopped")
	}

	// Check if taskId already exists
	if _, exists := s.taskMap[taskId]; exists {
		return fmt.Errorf("task with id %s already exists", taskId)
	}

	runAt := time.Now().Add(delay)
	task := &Task{
		id:      taskId,
		runAt:   runAt,
		payload: payload,
		handler: nil, // Default handler, can be extended
	}

	s.taskMap[taskId] = task
	heap.Push(&s.taskHeap, task)

	// Signal that a new task was added
	select {
	case s.nextTaskCh <- struct{}{}:
	default:
		// Channel already has a signal, no need to block
	}

	return nil
}

// Cancel cancels a pending task
func (s *Scheduler) Cancel(taskId string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.taskMap[taskId]
	if !exists {
		return false
	}

	// Remove from map
	delete(s.taskMap, taskId)

	// Note: We can't efficiently remove from heap, so we mark it as cancelled
	// The run loop will skip tasks that aren't in the map
	// For a more efficient implementation, we could add a cancelled flag to Task
	// and check it before execution

	return true
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if s.stopped || !s.started {
		s.mu.Unlock()
		return
	}
	s.stopped = true
	s.mu.Unlock()

	// Signal stop
	close(s.stopCh)

	// Wait for the run loop to finish
	s.wg.Wait()

	// Clear pending tasks
	s.mu.Lock()
	s.taskMap = make(map[string]*Task)
	s.taskHeap = make(TaskHeap, 0)
	s.mu.Unlock()
}

