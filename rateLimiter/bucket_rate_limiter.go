package main

import (
	"sync"
	"time"
)

// BucketRateLimiter - rate limiter using token bucket algo
type BucketRateLimiter struct {
	tokens int
	window time.Duration
	rMap sync.Map
	refillRate float64
}

type Bucket struct{
	rw sync.Mutex
	lastFilledAt int64
	tokensLeft int
}

func NewBucketRateLimiter(tokens int, window time.Duration) *BucketRateLimiter {
	return &BucketRateLimiter{
		tokens: tokens,
		window: window,
		refillRate : float64(tokens) / float64(window.Nanoseconds()),
	}
}

func (rl *BucketRateLimiter) Allow(id int) bool {
	currTime := time.Now().UnixNano()
	actual, loaded := rl.rMap.LoadOrStore(id, &Bucket{
		lastFilledAt: currTime,
		tokensLeft:   rl.tokens - 1,
	})
	if !loaded {
		return true
	}
	// retreive entry and lock
	entryVal := actual.(*Bucket)
	entryVal.rw.Lock()
	defer entryVal.rw.Unlock()

	// check we can accomodate requests or not
	
	// fill bucket corresponding to time elapsed
	elapsed := currTime - entryVal.lastFilledAt
	entryVal.tokensLeft += int(float64(elapsed) * rl.refillRate)
	if entryVal.tokensLeft > rl.tokens {
		entryVal.tokensLeft = rl.tokens
	}
	entryVal.lastFilledAt = currTime
	// no tokens remaining && cannot fill the bucket
	if entryVal.tokensLeft == 0 {
		return false
	}

	// use 1 token to accomodate current request
	entryVal.tokensLeft = entryVal.tokensLeft-1
	// update the entry in the map
	rl.rMap.Store(id, entryVal)
	return true
}