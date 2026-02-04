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


# Go Concurrency Practice — Extensions (Level 2 → Level 3)

This document extends the **basic goroutines + channels exercises** into
real-world concurrency patterns used in backend systems and interviews.

Focus areas:
- worker pools
- fan-in / fan-out
- cancellation
- rate limiting
- scheduling
- error propagation
- backpressure

The goal is **correctness > cleverness**.

---

## Level 2 — Concurrency Patterns

### Q2.1: Worker Pool (Bounded Concurrency)

**Concept:**  
Limit the number of goroutines processing jobs concurrently.

**Problem:**
- Jobs: integers `1..20`
- Start **3 worker goroutines**
- Each worker:
  - reads from `jobs` channel
  - sleeps for `300ms`
  - prints:  
    ```
    worker <id> processed job <jobId>
    ```
- Main:
  - sends all jobs
  - closes the jobs channel
  - waits for all workers to finish

**Constraints:**
- No `time.Sleep` in `main`
- Use `sync.WaitGroup`
- No global variables

**What this tests:**
- correct channel closing
- worker lifecycle management
- avoiding goroutine leaks

---

### Q2.2: Fan-out → Fan-in (Merge Channels)

**Concept:**  
Multiple producers → single consumer.

**Problem:**
- Producer A sends **even numbers** `0..10`
- Producer B sends **odd numbers** `1..9`
- Each producer has its **own channel**
- Merge both into **one output channel**
- Main ranges over output channel and prints values

**Hints:**
- Use `select`
- Close output channel **only after both producers finish**

**What this tests:**
- coordination correctness
- avoiding premature close
- select-based merging

---

### Q2.3: Cancellation-aware Worker Pool

**Concept:**  
Graceful shutdown using cancellation signals.

**Problem:**
- Extend Q2.1 worker pool
- Add a `done` channel
- Workers must:
  - stop immediately when `done` is closed
  - exit cleanly
- Main:
  - starts workers
  - sends jobs
  - cancels execution after `1 second`

**Rules:**
- No goroutine should block forever
- No job should be processed after cancellation

**What this tests:**
- cancellation patterns
- select usage
- goroutine lifecycle discipline

---

### Q2.4: Rate Limiter (Token Bucket – Simplified)

**Concept:**  
Control throughput using time.

**Problem:**
- Jobs: `1..10` arrive immediately
- Allow **only 1 job every 500ms**
- Use a `time.Ticker`
- Process job only when:
  - a token is available
  - a job exists

**Hint:**
```go
ticker := time.NewTicker(500 * time.Millisecond)
