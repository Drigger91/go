package main

type KeyValueStore[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, value V) (bool, error)
	Delete(key K)
}