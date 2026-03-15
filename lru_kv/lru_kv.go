package main

import (
	"fmt"
	"sync"
)

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