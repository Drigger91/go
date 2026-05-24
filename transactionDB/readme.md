# Machine Coding Round: TransactionDB

## Problem Statement

Build a **transactional key-value store** in Go — an in-memory database that supports basic CRUD operations along with nested transactions (BEGIN / COMMIT / ROLLBACK).

This is a **60–90 minute** machine coding exercise. Focus on correctness first, then clean design.

---

## Functional Requirements

### Basic Key-Value Operations

| Command | Description |
|---|---|
| `SET key value` | Store a value under the given key |
| `GET key` | Retrieve the value for a key; return an error/sentinel if not found |
| `DELETE key` | Remove a key from the store |
| `EXISTS key` | Return true/false whether the key exists |

### Transaction Commands

| Command | Description |
|---|---|
| `BEGIN` | Start a new transaction (transactions can be nested) |
| `COMMIT` | Persist all changes made in the current transaction to the outer scope |
| `ROLLBACK` | Discard all changes made in the current transaction |

### Count / Query Operations

| Command | Description |
|---|---|
| `COUNT value` | Return the number of keys that currently hold the given value |

---

## Behaviour Specification

### Outside a Transaction
- `GET`, `SET`, `DELETE`, `EXISTS`, `COUNT` operate directly on the main store.
- Calling `COMMIT` or `ROLLBACK` outside an active transaction should return a meaningful error.

### Inside a Transaction
- Changes made after `BEGIN` are **isolated** to that transaction — they are **not** visible to the outer scope until `COMMIT`.
- `ROLLBACK` discards all changes since the last `BEGIN` and returns to the state at the time of that `BEGIN`.
- `COMMIT` merges changes upward — into the enclosing transaction if nested, or into the main store if at the outermost level.

### Nested Transactions
```
BEGIN           # tx depth: 1
  SET a 1
  BEGIN         # tx depth: 2
    SET a 2
    ROLLBACK    # tx depth back to 1 → a = 1
  COMMIT        # tx depth: 0 → a = 1 committed to main store
```

- Each `BEGIN` pushes a new transaction scope onto a stack.
- Each `COMMIT` pops the top scope and merges its changes one level up.
- Each `ROLLBACK` pops the top scope and discards its changes entirely.

---

## Example Walkthrough

```
SET x 10
GET x           → 10

BEGIN
  SET x 20
  GET x         → 20       ← sees local change
ROLLBACK
GET x           → 10       ← rolled back to original

BEGIN
  SET x 30
  BEGIN
    SET x 40
    COMMIT              ← merges x=40 into outer tx
  GET x         → 40
COMMIT                  ← merges x=40 into main store
GET x           → 40

COUNT 40        → 1
EXISTS x        → true
DELETE x
EXISTS x        → false
COUNT 40        → 0
```

---

## Interface / API Design (Suggested)

You are free to choose your own design, but your solution must expose at minimum:

```go
type DB interface {
    Set(key, value string)
    Get(key string) (string, error)
    Delete(key string)
    Exists(key string) bool
    Count(value string) int

    Begin()
    Commit() error
    Rollback() error
}
```

You may add helper types, internal structs, or extend this interface as you see fit.

---

## Non-Functional Requirements

- **Language**: Go (no external packages beyond the standard library)
- **Concurrency**: Single-threaded is acceptable for this exercise; thread safety is a bonus
- **No persistence**: Everything lives in memory
- **Error handling**: Use idiomatic Go error returns — no panics for expected error states

---

## Deliverables

1. `db.go` — Core implementation
2. `db_test.go` — Unit tests covering the scenarios above (use `testing` package)
3. `main.go` (optional) — A simple REPL or CLI driver for manual testing

---

## Evaluation Criteria

| Area | What we look at |
|---|---|
| **Correctness** | All operations behave as specified, especially nested transactions |
| **Data structures** | Appropriate choice for the transaction stack and scope isolation |
| **Code clarity** | Idiomatic Go, clean naming, small focused functions |
| **Error handling** | Proper use of `error` returns; no silent failures |
| **Test coverage** | Edge cases: nested rollback, commit with no begin, delete inside tx, etc. |
| **Extensibility** | How hard would it be to add `SAVEPOINT`? TTL? Logging? |

---

## Constraints & Rules

- Do **not** use any database library (e.g., BoltDB, SQLite, BadgerDB).
- Do **not** look up solutions online during the exercise.
- You may use the Go standard library freely (`fmt`, `errors`, `strings`, etc.).
- Ask clarifying questions **before** you start coding if anything is ambiguous.

---

## Hints (read only if stuck)

<details>
<summary>Hint 1 — Data structure</summary>

Think of the transaction stack as a `[]map[string]string`. Each `BEGIN` pushes a new empty map. Operations write to the top-most map. `COMMIT` merges the top map into the one below. `ROLLBACK` pops and discards.

</details>

<details>
<summary>Hint 2 — Deletions inside a transaction</summary>

You need to distinguish between "key does not exist in this layer" and "key was explicitly deleted in this layer." Consider using a sentinel value or a separate deleted-keys set per transaction frame.

</details>

<details>
<summary>Hint 3 — COUNT efficiency</summary>

A naive implementation scans all keys on every `COUNT` call (O(n)). For a bonus, maintain a reverse index `map[string]int` (value → count) that you update on every `SET` and `DELETE`.

</details>

---

## Stretch Goals (if time permits)

- [ ] Thread-safe implementation using `sync.RWMutex`
- [ ] `SAVEPOINT name` / `ROLLBACK TO name` — named savepoints within a transaction
- [ ] TTL support: `SET key value ttl_seconds`
- [ ] A simple line-by-line REPL (`> SET x 1`, `> GET x`)
- [ ] Undo log / write-ahead log printed to stdout for observability