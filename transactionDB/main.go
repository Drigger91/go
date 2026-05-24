package main

import "fmt"

func main() {
	db := NewBasicTransactionDB()

	fmt.Println("=== Basic Operations ===")
	db.Set("x", 10)
	db.Set("color", "blue")
	fmt.Println("x:", db.Get("x"))
	fmt.Println("color:", db.Get("color"))
	fmt.Println("exists x:", db.Exists("x"))
	fmt.Println("exists y:", db.Exists("y"))

	fmt.Println("\n=== Delete ===")
	err := db.Delete("color")
	fmt.Println("delete color:", err)
	fmt.Println("color after delete:", db.Get("color"))
	err = db.Delete("color")
	fmt.Println("delete color again:", err)

	fmt.Println("\n=== Single Transaction: Rollback ===")
	db.Set("x", 10)
	db.Begin()
	db.Set("x", 99)
	fmt.Println("x inside tx:", db.Get("x"))
	db.Rollback()
	fmt.Println("x after rollback:", db.Get("x"))

	fmt.Println("\n=== Single Transaction: Commit ===")
	db.Begin()
	db.Set("x", 42)
	fmt.Println("x inside tx:", db.Get("x"))
	db.Commit()
	fmt.Println("x after commit:", db.Get("x"))

	fmt.Println("\n=== Nested Transactions ===")
	db.Set("x", 10)
	db.Begin()
	db.Set("x", 20)
	db.Begin()
	db.Set("x", 30)
	fmt.Println("x at depth 2:", db.Get("x"))
	db.Rollback()
	fmt.Println("x after inner rollback:", db.Get("x"))
	db.Commit()
	fmt.Println("x after outer commit:", db.Get("x"))

	fmt.Println("\n=== Delete Inside Transaction ===")
	db.Set("y", "hello")
	db.Begin()
	db.Delete("y")
	fmt.Println("y inside tx (should be nil):", db.Get("y"))
	fmt.Println("exists y inside tx:", db.Exists("y"))
	db.Rollback()
	fmt.Println("y after rollback (should be hello):", db.Get("y"))

	fmt.Println("\n=== Commit/Rollback with no active tx ===")
	fmt.Println("commit:", db.Commit())
	fmt.Println("rollback:", db.Rollback())

}
