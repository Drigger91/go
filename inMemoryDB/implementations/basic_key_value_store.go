package implementations

import (
	"reflect"
	"sync"
)

type BasicKeyValueStore[K comparable] struct {
	store map[K]any // this will contain key - entry mappings
	rw sync.RWMutex
}
func NewBasicKeyValueStore[K comparable]() BasicKeyValueStore[K] {
	return BasicKeyValueStore[K]{
		store: make(map[K]any),
	}
}
func (kv *BasicKeyValueStore[K]) Set(key K, newValue any) (bool, error) {
	kv.rw.Lock()
	defer kv.rw.Unlock()
	if _, exists := kv.store[key]; !exists {
		kv.store[key] = newValue
		return true, nil	
	}
	// already exist
	val := kv.store[key]
	if reflect.TypeOf(val) != reflect.TypeOf(newValue) {
		panic("Types are not matching for existing key")
	}
	kv.store[key] = newValue
	return true, nil
}

func (kv *BasicKeyValueStore[K]) Get(key K) (any, bool) {
	kv.rw.RLock()
	defer kv.rw.RUnlock()
	value, exists := kv.store[key]
	return value, exists
}
func (db *BasicKeyValueStore[K]) Delete(key K) {
	db.rw.Lock()
	defer db.rw.Unlock()
	delete(db.store, key)
}