# Go Concurrency + DS Practice (Small → Medium → LLD)

This README contains a step-by-step practice ladder to build strong fundamentals in:
- goroutines
- channels
- coordination patterns
- timers + select
- worker pools
- core data structures (heap, LRU)
- building blocks for machine coding / LLD

Goal: **start small**, build confidence, then scale into bigger systems like a scheduler.

---

## Level 0 — Goroutines + Channels Basics

### Q0.1: Hello from Goroutine
**Problem:**  
Start a goroutine that prints `"hello from goroutine"` and main prints `"hello from main"`.

**Expected:** Both lines print (order not guaranteed).

---

### Q0.2: Unbuffered Channel Send/Receive
**Problem:**  
Create an unbuffered channel `chan string`.
- Start a goroutine that sends `"hello"` after sleeping 1 second.
- Main receives and prints it.

**Goal:** Understand blocking behavior.

---

### Q0.3: Buffered Channel Behavior
**Problem:**  
Create a buffered channel `chan int` with capacity `2`.
- Send 2 values without receiving
- Try sending the 3rd and observe it blocks
- Add a receiver and verify it unblocks

---

### Q0.4: Close + Range
**Problem:**  
Producer goroutine sends integers `1..5` into a channel and then closes it.
Main goroutine ranges over channel and prints all numbers.

**Goal:** Learn `close(ch)` and `for x := range ch`.

---

---

## Level 1 — Coordination Patterns (1–2 hrs)

### Q1.1: Signal Once (Channel based)
**Problem:**  
- Main starts a goroutine that waits on a signal channel `startCh`
- Main sleeps 500ms then sends a signal
- Goroutine prints `"started"` and exits
- Main prints `"done"`

**Constraint:** No WaitGroup.

---

### Q1.2: Done Channel Cancellation
**Problem:**  
- Start a goroutine that prints `"working..."` every 300ms
- Main waits 2 seconds, then closes `done` channel
- Worker should stop immediately and print `"stopped"`

**Goal:** Learn cancellation patterns without killing goroutines forcefully.

---

### Q1.3: Timeout using `select`
**Problem:**  
Write a function:

```go
func waitForMessage(ch <-chan string, timeout time.Duration) (string, bool)
It waits for a message on ch

If message arrives before timeout → return (msg, true)

If timeout happens first → return ("", false)

Use:
- select
- time.After(timeout) 