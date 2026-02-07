package main

import (
	"fmt"

	"github.com/in/Drigger91/inMemDb/implementations"
)


func main() {
	fmt.Println("In mem db")

	store := implementations.NewBasicKeyValueStore[any]()

	store.Set("1", "one")
	store.Set(1, "int one")

	fmt.Println(store.Get("1"))
	fmt.Println(store.Get(1))

	// _, err := store.Set("1", 1)
	// fmt.Println("Err:", err)

	store.Delete("1")

	store.Set("1", [3]string{"1", "2", "3"})
	fmt.Println(store.Get("1"))
	store.Set("1", [3]int{1})
}