package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Scheduled Queue Example")
	sch := NewScheduler()
	sch.Start()
	err := sch.Schedule("t1", 2*time.Second, "email:user=1")
	if err != nil {
		fmt.Printf("Error scheduling t1: %v\n", err)
	}
	err = sch.Schedule("t2", 1*time.Second, "push:user=2")
	if err != nil {
		fmt.Printf("Error scheduling t2: %v\n", err)
	}
	err = sch.Schedule("t3", 3*time.Second, "sms:user=3")
	if err != nil {
		fmt.Printf("Error scheduling t3: %v\n", err)
	}
	time.Sleep(1500 * time.Millisecond)
	cancelled := sch.Cancel("t3")
	if cancelled {
		fmt.Println("Cancelled task t3")
	} else {
		fmt.Println("Task t3 not found or already executed")
	}
	time.Sleep(4 * time.Second)
	sch.Stop()
	fmt.Println("Scheduler stopped")
}