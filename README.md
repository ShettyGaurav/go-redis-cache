## **GoCache** – A Redis-like Cache System in Go

  GoCache is a lightweight Redis-inspired cache system written entirely in Go.
It supports multiple data types, file-based persistence(For Recovery)

Data Types Supported

String

List (LPUSH, LPOP, LRANGE)

Set (planned)

Hash (planned)

Persistence

Cache entries are stored in a file.
On restart, non-expired data is loaded back into memory.

Core Commands

SET – Store a string value with optional expiration.

GET – Retrieve string value.

DELETE – Remove a key.

EXISTS – Check if a key exists.

GETTYPE – Get the type of a key.

LPUSH – Push an element to the head of a list.

LPOP – Pop an element from the head of a list.

LRANGE – Get a slice of elements from a list.

## 🛠 Guide

Clone the repository and build:

```bash
git clone https://github.com/ShettyGaurav/go-redis-cache.git
go mod tidy
cd go-redis-cache
go run main.go
```
## 🛠 Under Progress
-Building CLI 

-Client-Server for cache Integration
