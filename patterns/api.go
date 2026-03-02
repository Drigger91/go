package main

import (
	"fmt"
	"time"
)


func fetchUser() string {
	time.Sleep(120 * time.Millisecond)
	return "user-data"
}

func fetchOrders() string {
	time.Sleep(150 * time.Millisecond)
	return "orders-data"
}

func fetchPayments() string {
	time.Sleep(100 * time.Millisecond)
	return "payments-data"
}

type Transaction struct {
	ID int
}

func fetchPage(offset, limit int) []Transaction {
	time.Sleep(120 * time.Millisecond) // simulate network latency

	var result []Transaction
	for i := offset; i < offset+limit; i++ {
		result = append(result, Transaction{ID: i})
	}

	fmt.Println("Fetched page:", offset)
	return result
}