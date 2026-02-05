# Machine Coding Problem — Rate Limiter (Go)

## Overview

You are building a backend service that must protect itself from abuse and overload.
To achieve this, you need to design and implement an **in-memory rate limiter**.

This problem tests:
- concurrency correctness
- time-based logic
- API design
- production-minded trade-offs

---

## Problem Statement

Design and implement a **thread-safe rate limiter** that limits how frequently a client can perform an action.

Each client is identified by a `key` (e.g. userId, IP address, API key).

---

## Core Requirements (Mandatory)

### 1. Public API

Expose the following interface:

```go
type RateLimiter interface {
    Allow(key string) bool
}
```

- Rate Limiting Rule

Allow N requests per T duration
Example:

5 requests per 1 second

burst up to N is allowed

After the time window resets, requests should be allowed again

-  Concurrency Requirements

Must be safe for concurrent access

Multiple goroutines may call Allow() simultaneously

No data races

No goroutine leaks

- Performance Requirements

Allow() should be O(1) or amortized O(1)

Avoid unnecessary allocations

Avoid per-request goroutines