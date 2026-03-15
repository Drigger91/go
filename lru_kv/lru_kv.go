package main

import (
	"fmt"
	"sync"
)
// important: the row level locking is technically not possible due to Get() getting called with non existent keys (good finding this yourself)
// secondly if the global mutex becomes the bottleneck (eventually it will for millions of users), we should still not try and implement row
// level locks as it is fundamentally flawed for this usecase (memory explosion of LoadOrStore()). We should move towards client side sharding kinda logic
// it will involve one more layering, something like: KeyValueStore -> map[string]*LruShards (lrushards being the current implementation),
//  and a func to calculate hash and pointing towards that shard. always remember not to complicate things with stuff like row level locks :)
type KeyValueStoreLRU struct {
	capacity int
	// to store the values
	dataStore map[string]*LruNode

	// lru Node for eviction policy
	lruNode *LRU

	// global rwmutex
	rw *sync.RWMutex
}

func NewKeyValueStore(capacity int) *KeyValueStoreLRU{
	return &KeyValueStoreLRU{
		capacity: capacity,
		dataStore: make(map[string]*LruNode),
		lruNode: NewLRU(),
		rw: &sync.RWMutex{},
	}
}


func (kv *KeyValueStoreLRU) Get(key string) (any, error) {
	kv.rw.Lock()
	defer kv.rw.Unlock()
	// check if the key exist in the datastore or not
	nodeValue, exists := kv.dataStore[key]

	if !exists {
		return nil, fmt.Errorf("Key Not found")
	}
	// update LRU here
	kv.lruNode.Update(nodeValue)

	return nodeValue.val, nil
}

func (kv *KeyValueStoreLRU) Put(key string, value any) bool {
	kv.rw.Lock()
	defer kv.rw.Unlock()

	// check whether this is a new key for size constraint
	_, exists := kv.dataStore[key]

	if !exists {
		// check if LRU needs to evicted
		currSize := len(kv.dataStore)

		if currSize == kv.capacity {
			// kv store already at capacity, lru needs to be evicted
			keyToBeDeleted := kv.lruNode.Evict()
			delete(kv.dataStore, keyToBeDeleted)
		}
		
	}

	// update the value in the datastore
	if !exists {
		kv.dataStore[key] = NewNode(key, value)
	} else {
		node := kv.dataStore[key]
		node.val = value
	}
	
	node := kv.dataStore[key]
	// update the lru here
	kv.lruNode.Update(node)
	return true
}  