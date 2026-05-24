package main

type KVStore interface {
	Get(key string) any
	Exists(key string) bool
	Set(key string, value any)
	Delete(key string) error
}

type Transactor interface {
	Begin()
	Commit() error
	Rollback() error
}

// db will commit, transaction not as a separate object
type SimpleTransactionDB interface {
	KVStore
	Transactor
}

const (
	SET      = "SET"
	GET      = "GET"
	BEGIN    = "BEGIN"
	ROLLBACK = "ROLLBACK"
	COMMIT   = "COMMIT"
	EXISTS   = "EXISTS"
	DELETE   = "DELETE"
)
