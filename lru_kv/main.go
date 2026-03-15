package main

import "fmt"

func main() {
	fmt.Println("LRU KV")
	kvStore := NewKeyValueStore(2)

	kvStore.Put("key1", "value1")
	kvStore.Put("key2", "value2")
	kvStore.Put("key3", "value3")
	
}