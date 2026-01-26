# Scheduled Task Queue (Delayed Job Runner) — Machine Coding Problem

## Overview
Build a **Scheduled Queue / Delayed Job Runner** that can accept tasks scheduled for the future and execute them **at (or after)** their scheduled time.

This problem tests:
- Priority ordering by time
- Concurrency + coordination
- Clean API design
- Graceful shutdown

---

## Requirements

### Task Model
Each task has:
- `id` (string)
- `runAt` (time.Time) — when it should execute
- `payload` (string)
- `handler` (function) — executed when the task runs

You can represent this as a struct.

---

## APIs to Implement

### 1) `NewScheduler() *Scheduler`
Creates a scheduler instance.

---

### 2) `Start()`
Starts the scheduler background loop.

- The scheduler should run tasks in a separate goroutine.
- Calling `Start()` multiple times should **not** create multiple loops.

---

### 3) `Schedule(taskId string, delay time.Duration, payload string) error`
Schedules a task to run **after** `delay`.

Rules:
- `runAt = time.Now() + delay`
- Tasks must execute in increasing order of `runAt`
- If 2 tasks have the same `runAt`, execute in FIFO order (first scheduled first)
- `taskId` must be unique
  - If a task with same `taskId` already exists (pending), return error

---

### 4) `Cancel(taskId string) bool`
Cancels a pending task.

- Returns `true` if task existed and was cancelled
- Returns `false` if task doesn't exist (or already executed)

---

### 5) `Stop()`
Gracefully stops the scheduler.

Rules:
- No new tasks should be accepted after stop
- Already executing tasks can finish
- Pending tasks should NOT run after stop
- Stop should not deadlock / block forever

---

## Execution Rules
When a task executes:
- Print/log something like:  
  `Executed task=<id> payload=<payload> at=<time>`
- Execute tasks **at or after** their scheduled time (not before)
- If the scheduler is idle, it should not consume CPU (no busy wait)

---

## Concurrency Expectations
- `Schedule()` can be called concurrently from multiple goroutines
- `Cancel()` can be called concurrently
- Scheduler must remain thread-safe and correct

---

## Error Handling
- If `handler(payload)` returns an error, log it but continue running further tasks
- Scheduler should not crash if one task fails

---

## Suggested Data Structures
You will likely need:
- A **min-heap / priority queue** ordered by `runAt`
- A map for `taskId -> task` lookup (for uniqueness + cancel)

---

## Bonus (Optional)
If you want extra challenge:
1. Add `ScheduleAt(taskId string, runAt time.Time, payload string) error`
2. Add retries with backoff:
   - If a task fails, retry up to `maxRetries`
3. Metrics:
   - tasks scheduled
   - tasks executed
   - tasks cancelled
   - execution delay (`actualRunTime - runAt`)

---

## Example Usage (Expected Behavior)

```go
sch := NewScheduler()
sch.Start()

sch.Schedule("t1", 2*time.Second, "email:user=1")
sch.Schedule("t2", 1*time.Second, "push:user=2")
sch.Schedule("t3", 3*time.Second, "sms:user=3")

time.Sleep(1500 * time.Millisecond)
sch.Cancel("t3")

time.Sleep(4 * time.Second)
sch.Stop()
