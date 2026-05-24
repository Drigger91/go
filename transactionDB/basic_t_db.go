package main

import "fmt"

// IMP: this is necessary to maintain the states (basically like a delimiter)
const deleted = "/_/__DELETED__/_/"

type state struct {
	data map[string]any
}

type BasicTransactionDB struct {
	dataStore    map[string]any
	transactions []*state
}

func NewBasicTransactionDB() SimpleTransactionDB {
	return &BasicTransactionDB{
		dataStore:    make(map[string]any),
		transactions: make([]*state, 0),
	}
}

func (db *BasicTransactionDB) activeState() *state {
	if len(db.transactions) == 0 {
		return nil
	}
	return db.transactions[len(db.transactions)-1]
}

func (db *BasicTransactionDB) Get(key string) any {
	for i := len(db.transactions) - 1; i >= 0; i-- {
		val, exists := db.transactions[i].data[key]
		if exists {
			if val == deleted {
				return nil
			}
			return val
		}
	}

	val, exists := db.dataStore[key]
	if !exists {
		return nil
	}
	return val
}

func (db *BasicTransactionDB) Set(key string, value any) {
	if f := db.activeState(); f != nil {
		f.data[key] = value
		return
	}
	db.dataStore[key] = value
}

func (db *BasicTransactionDB) Exists(key string) bool {
	return db.Get(key) != nil
}

func (db *BasicTransactionDB) Delete(key string) error {
	if !db.Exists(key) {
		return fmt.Errorf("the key %s does not exist", key)
	}

	if f := db.activeState(); f != nil {
		f.data[key] = deleted
		return nil
	}
	delete(db.dataStore, key)
	return nil
}

func (db *BasicTransactionDB) Begin() {
	db.transactions = append(db.transactions, &state{data: make(map[string]any)})
}

func (db *BasicTransactionDB) Commit() error {
	if len(db.transactions) == 0 {
		return fmt.Errorf("no active transaction to commit")
	}

	top := db.activeState()
	db.transactions = db.transactions[:len(db.transactions)-1]

	var target map[string]any
	if len(db.transactions) > 0 {
		target = db.activeState().data
	} else {
		target = db.dataStore
	}

	for k, v := range top.data {
		if v == deleted {
			delete(target, k)
		} else {
			target[k] = v
		}
	}
	return nil
}

func (db *BasicTransactionDB) Rollback() error {
	if len(db.transactions) == 0 {
		return fmt.Errorf("no active transaction to rollback")
	}
	db.transactions = db.transactions[:len(db.transactions)-1]
	return nil
}

func Rollback() error {
	return nil
}