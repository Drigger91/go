# LRU Cache – Machine Coding Problem

## Problem Statement

Design and implement a **Least Recently Used (LRU) Cache**.

The cache should support constant time operations for retrieving and inserting key-value pairs. When the cache reaches its maximum capacity, it should evict the **least recently used** item before inserting a new one.

Your implementation should be **efficient, extensible, and thread-safe (optional extension)**.

---

# Functional Requirements

## 1. Initialize Cache

Create a cache with a fixed capacity.

```
NewLRUCache(capacity int)
```

Example:

```
cache := NewLRUCache(3)
```

---

## 2. Get Value

```
Get(key)
```

Behavior:

- If key exists → return value
- Mark the key as **most recently used**
- If key does not exist → return `-1` or equivalent

Example:

```
cache.Get(2)
```

---

## 3. Put Key-Value Pair

```
Put(key, value)
```

Behavior:

- Insert key-value pair if it does not exist
- Update value if key already exists
- Mark the key as **most recently used**

If cache capacity is exceeded:

- Remove the **least recently used** item.

---

# Example Flow

```
capacity = 2

Put(1,1)
Put(2,2)

Get(1) -> returns 1

Put(3,3)
→ evicts key 2

Get(2) -> returns -1

Put(4,4)
→ evicts key 1

Get(1) -> -1
Get(3) -> 3
Get(4) -> 4
```

---

# Constraints

You may assume:

```
1 <= capacity <= 10^5
number of operations <= 10^6
```

Required complexity:

```
Get → O(1)
Put → O(1)
```

---

# Expected Design

The typical solution uses:

```
HashMap + Doubly Linked List
```

### HashMap

```
key → node reference
```

Provides O(1) lookup.

---

### Doubly Linked List

Maintains usage order.

```
Head → Most Recently Used
Tail → Least Recently Used
```

Example structure:

```
[HEAD]
   |
[3] <-> [5] <-> [7]
               |
            [TAIL]
```

- New access → move node to **head**
- Eviction → remove node from **tail**

---

# Suggested Data Structures

```
LRUCache
 ├── capacity
 ├── hashmap[key] → node
 ├── head
 └── tail
```

Node structure:

```
Node
 ├── key
 ├── value
 ├── prev
 └── next
```

---

# Operations Overview

### Get

```
1. lookup node in hashmap
2. move node to head
3. return value
```

---

### Put

```
1. check if key exists
2. if exists → update value + move to head
3. if new → create node
4. add node to head
5. if capacity exceeded → remove tail node
```

---

# Edge Cases

Your implementation should handle:

```
capacity = 1
duplicate keys
Get on missing key
multiple updates to same key
```

---

# Optional Extensions (SDE-2 Level)

If time permits, extend your design.

### 1. Thread-Safe Cache

Ensure concurrent calls to:

```
Get()
Put()
```

do not corrupt the cache state.

Possible approaches:

```
Mutex
RWMutex
```

---

### 2. Time-To-Live (TTL)

Support expiration:

```
Put(key, value, ttl)
```

Expired entries should be removed automatically.

---

### 3. Metrics

Expose cache metrics:

```
CacheHits
CacheMisses
Evictions
```

---

# Evaluation Criteria

Interviewers usually judge:

### Correctness

```
O(1) Get
O(1) Put
correct eviction
```

---

### Code Structure

```
clean abstractions
readable code
proper data ownership
```

---

### Edge Case Handling

```
capacity overflow
duplicate keys
invalid lookups
```

---

# Expected Interview Time

Typical machine coding round:

```
45–60 minutes
```

Breakdown:

```
10 min → design discussion
30 min → coding
10 min → extensions
```

---

# Suggested Development Order

To avoid overengineering:

```
1. Basic LRU with hashmap + doubly linked list
2. Correct eviction logic
3. Edge cases
4. Thread safety (optional)
5. Extensions
```