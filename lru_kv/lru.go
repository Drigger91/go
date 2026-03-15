package main

import (
	"fmt"
	"math"
)

type LruNode struct {
	val any
	key string
	next *LruNode
	prev *LruNode
}

func NewNode(key string, val any) *LruNode {
	return &LruNode{
		val: val,
		key: key,
		next: nil,
		prev: nil,
	}
}

type LRU struct {
	head *LruNode // least recently used value will be next to head
	tail *LruNode // most frequent will be prev to tail
}
func NewLRU() *LRU {
	lru := &LRU{
		head: NewNode("HEAD", math.MaxInt),
		tail: NewNode("TAIL" , math.MaxInt),
	}

	headNode := lru.head
	tailNode := lru.tail
	headNode.next = tailNode
	tailNode.prev = headNode

	return lru
}

// Update - takes node as param, makes node the most frequently used node
func (lru *LRU) Update(node *LruNode) {
	prevNode := node.prev
	nextNode := node.next

	// sever the link
	if prevNode != nil {
		prevNode.next = nextNode
	}
	if nextNode != nil {
		nextNode.prev = prevNode
	}

	prevMFUnode := lru.tail.prev

	// create the link btw prev MFU and new MFU
	prevMFUnode.next = node
	node.prev = prevMFUnode

	// create the link btw new MFU and tail
	node.next = lru.tail
	lru.tail.prev = node

}

func (lru *LRU) Evict() string {

	lruNode := lru.head.next

	if lru.head.next == lru.tail {
		return "" 
	}

	lruNode.prev.next = lruNode.next
	lruNode.next.prev = lruNode.prev

	// free up (to check whether this is correct or not)
	key := lruNode.key
	fmt.Println("Evicting lru key:", key)
	return key
}


