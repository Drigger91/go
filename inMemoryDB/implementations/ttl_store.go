package implementations

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)
var inf time.Time
type KeyValueEntry struct {
	value any
	ttl time.Time
}
type TtlIndex[K comparable] struct {
	keys   []K
	cursor int 
	mu     sync.Mutex
}
type KeyValueStore[K comparable] struct {
	store map[K]KeyValueEntry // this will contain key - entry mappings
	rw sync.RWMutex
	ttlIndex TtlIndex[K]
	stopCh   chan struct{}   // signals shutdown
	wg       sync.WaitGroup 
	closed   bool
}
func NewKeyValueStore[K comparable]() *KeyValueStore[K] {
	// register the gc here.
	stopCh := make(chan struct{}, 1)
	kv := &KeyValueStore[K]{
		store: make(map[K]KeyValueEntry),
		stopCh: stopCh,
	}
	// fire up in different routine
	go kv.ttlCleaner()
	return kv
}

func (kv *KeyValueStore[K]) ttlCleaner() {
	defer kv.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			kv.cleanupTtlKeys(500) // bounded work
		case <-kv.stopCh:
			return
		}
	}
}

func (kv *KeyValueStore[K]) cleanupTtlKeys(budget int) {
	fmt.Println("GC invoked")
	kv.ttlIndex.mu.Lock()
	if len(kv.ttlIndex.keys) == 0 {
		kv.ttlIndex.mu.Unlock()
		return
	}

	start := kv.ttlIndex.cursor
	end := min(start + budget, len(kv.ttlIndex.keys))

	// slice of candidates for this cycle
	candidates := kv.ttlIndex.keys[start:end]
	kv.ttlIndex.cursor = end

	// reset cursor if we reached the end
	if kv.ttlIndex.cursor >= len(kv.ttlIndex.keys) {
		kv.ttlIndex.cursor = 0
	}
	kv.ttlIndex.mu.Unlock()

	now := time.Now()
	var compacted []K
	compactedNeeded := false

	kv.rw.Lock()
	for _, key := range candidates {
		entry, exists := kv.store[key]
		if !exists {
			compactedNeeded = true
			continue
		}

		if entry.ttl.IsZero() {
			compactedNeeded = true
			continue
		}

		if now.After(entry.ttl) {
			fmt.Println("deleting key", key)
			delete(kv.store, key)
			compactedNeeded = true
			continue
		}

		compacted = append(compacted, key)
	}
	kv.rw.Unlock()

	// Opportunistic compaction
	if compactedNeeded {
		kv.ttlIndex.mu.Lock()
		newSlice := kv.ttlIndex.keys[:start] // till start
		newSlice = append(newSlice, compacted...) // filtered compacted list
		newSlice = append(newSlice, kv.ttlIndex.keys[end:]...) // remaining from end
		kv.ttlIndex.keys = newSlice
		kv.ttlIndex.mu.Unlock()
	}
}
func (kv *KeyValueStore[K]) registerKeyForTtlIndex(key K) {
	kv.ttlIndex.mu.Lock()
	defer kv.ttlIndex.mu.Unlock()
	kv.ttlIndex.keys = append(kv.ttlIndex.keys, key)
}
func (kv *KeyValueStore[K]) SetEx(key K, newValue any, ttl time.Duration) (bool, error) {
	kv.rw.Lock()
	defer kv.rw.Unlock()
	entryVal := KeyValueEntry{
		value: newValue,
	}
	if ttl > 0 {
		entryVal.ttl = time.Now().Add(ttl)
		kv.registerKeyForTtlIndex(key)
	}
	if _, exists := kv.store[key]; !exists {
		kv.store[key] = entryVal
		return true, nil	
	}
	// already exist
	val := kv.store[key].value
	if reflect.TypeOf(val) != reflect.TypeOf(newValue) {
		return false, fmt.Errorf("Types are not matching for existing key")
	}

	kv.store[key] = entryVal
	return true, nil

}
func (kv *KeyValueStore[K]) Set(key K, newValue any) (bool, error) {
	return kv.SetEx(key, newValue, 0)
}

func checkForExpiry(value KeyValueEntry) bool {
	return !value.ttl.IsZero() && time.Now().After(value.ttl)
}

func (kv *KeyValueStore[K]) Get(key K) (any, bool) {
	kv.rw.RLock()
	defer kv.rw.RUnlock()
	value, exists := kv.store[key]
	// check if ttl is expired
	if exists && checkForExpiry(value) {
		return BasicKeyValueEntry{}, false
	}
	return value, exists
}
func (db *KeyValueStore[K]) Delete(key K) {
	db.rw.Lock()
	defer db.rw.Unlock()
	delete(db.store, key)
}

// Cleanup - deletes everything and clears the existing store
func (kv *KeyValueStore[K]) Cleanup() {
	kv.rw.Lock()
	if kv.closed {
		kv.rw.Unlock()
		return
	}
	kv.closed = true
	kv.rw.Unlock()

	// stop background goroutines
	close(kv.stopCh)

	// wait for cleaner to exit
	kv.wg.Wait()

	// clear all data safely
	kv.rw.Lock()
	clear(kv.store)
	kv.rw.Unlock()
}