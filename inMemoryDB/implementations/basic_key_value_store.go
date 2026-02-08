package implementations

import (
	"reflect"
	"sync"
	"time"
)

type BasicKeyValueEntry struct {
	value any
	ttl int64
}
type BasicKeyValueStore[K comparable] struct {
	store map[K]BasicKeyValueEntry // this will contain key - entry mappings
	rw sync.RWMutex
}
func NewBasicKeyValueStore[K comparable]() *BasicKeyValueStore[K] {
	return &BasicKeyValueStore[K]{
		store: make(map[K]BasicKeyValueEntry),
	}
}

func (kv *BasicKeyValueStore[K]) SetEx(key K, newValue any, ttl time.Duration) (bool, error) {
	kv.rw.Lock()
	defer kv.rw.Unlock()
	entryVal := BasicKeyValueEntry{
		value: newValue,
		ttl: time.Now().UnixNano() + ttl.Nanoseconds(),
	}
	if _, exists := kv.store[key]; !exists {
		kv.store[key] = entryVal
		return true, nil	
	}
	// already exist
	val := kv.store[key].value
	if reflect.TypeOf(val) != reflect.TypeOf(newValue) {
		panic("Types are not matching for existing key")
	}

	kv.store[key] = entryVal
	return true, nil

}
func (kv *BasicKeyValueStore[K]) Set(key K, newValue any) (bool, error) {
	// setting 5 years as infinity for now
	return kv.SetEx(key, newValue, 5 * time.Duration(time.Now().Year()))
}

func (kv *BasicKeyValueStore[K]) Get(key K) (any, bool) {
	kv.rw.RLock()
	defer kv.rw.RUnlock()
	value, exists := kv.store[key]
	// check if ttl is expired
	if exists && time.Now().UnixNano() > value.ttl {
		return BasicKeyValueEntry{}, false
	}
	return value, exists
}
func (db *BasicKeyValueStore[K]) Delete(key K) {
	db.rw.Lock()
	defer db.rw.Unlock()
	delete(db.store, key)
}