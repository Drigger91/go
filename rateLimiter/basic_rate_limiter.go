package main

import (
	"sync"
	"time"
)

type BasicRateLimiter struct {
	rMap *sync.Map
	tokens int
	window time.Duration
}

type Entry struct {
	rw sync.RWMutex
	Timestamps []int64
}

func NewBasicRateLimiter(tokens int, duration time.Duration) *BasicRateLimiter {
	return &BasicRateLimiter{
		rMap: new(sync.Map),
		tokens: tokens,
		window: duration,
	}
}


func (entry *Entry) binarySearch(k int64) int {
	s := 0
	e := len(entry.Timestamps)
	arr := entry.Timestamps

	if arr[0] >= k {
		return -1
	}
	ans := e
	for (s < e) {
		mid := s + (e-s)/2;
		if arr[mid] > k {
			e = mid-1
		} else {
			ans = min(ans, mid)
			s = mid+1
		}
	}
	return ans
}

func createNewEntry(id int, rMap *sync.Map, ts int64) {
	var timeStamps []int64
	timeStamps = append(timeStamps, ts)
	rMap.Store(id, &Entry{
		Timestamps: timeStamps,
	})
}

// Allow - based on userID, this method will decide whether to allow requests or not
func (rl *BasicRateLimiter) Allow(id int) bool {
	currTimeStamp := time.Now().UnixNano()
	if _, ok := rl.rMap.Load(id); !ok {
		createNewEntry(id, rl.rMap, currTimeStamp)
		return true 
	}
	// entry exists in map
	val, _ := rl.rMap.Load(id)
	// lock the mutex

	// get timestamps for the existing
	entry := val.(*Entry)

	// lock the mutex, and unlock after completion
	entry.rw.Lock()
	defer entry.rw.Unlock()
	timestamps := entry.Timestamps

	if len(timestamps) < rl.tokens {
		updateRmap(rl.rMap, currTimeStamp, timestamps, id)
		return true
	}
	// threshold
	k := time.Now().UnixNano() - rl.window.Nanoseconds()
	// fmt.Println("K" , k)
	// fmt.Println("TS", timestamps)
	upperBound := entry.binarySearch(k)

	if upperBound == -1 {
		// not possible to accomodate the request
		return false
	}

	// otherwise safedelete the tokens and add the current ones
	cut := min(upperBound + 1, len(timestamps))
	timestamps = timestamps[cut:]

	if len(timestamps) == rl.tokens {
		return false
	}
	updateRmap(rl.rMap, currTimeStamp, timestamps, id)
	return true
}

func updateRmap(rMap *sync.Map, newTimeStamp int64, timestamps []int64, id int) {
	timestamps = append(timestamps, newTimeStamp)
	rMap.Swap(id, &Entry{Timestamps: timestamps})
}