# search

A high-performance in-memory user matching engine written in Go.  
This package is designed to support matchmaking use cases such as dating apps, social platforms, or recommendation systems.  

It provides fast `Create`, `Search`, and `Delete` operations with zero allocations during search.

---

## Features

- Store user profiles with:
  - Age preferences
  - Sex preferences (`0` = male, `1` = female, `2` = any)
  - Language
  - Interests
  - Priority flag (boosts match score)
- Match users by multiple conditions:
  - Language
  - Age range compatibility
  - Sex compatibility
  - Shared interests
  - Priority
- Thread-safe operations
- Optimized for speed with minimal memory allocations

---

## Installation

```bash
go get github.com/elum-utils/search
````

---

## Usage

```go
package main

import (
	"fmt"
	"github.com/elum-utils/search"
)

func main() {
	// Create a new user
	search.Create(
		1,          // UserID
		"en",       // Language
		18, 30,     // Preferred partner age range
		1,          // Preferred partner sex (1 = female)
		25,         // My age
		0,          // My sex (0 = male)
		false,      // Priority flag
		"music", "art", "games", // Interests
	)

	// Another user searching for a match
	match := search.Search(
		2,          // My ID
		"en",       // My language
		20, 28,     // My preferred partner age range
		0,          // My preferred partner sex
		24,         // My age
		1,          // My sex
		"music", "art", // My interests
	)

	if match != nil {
		fmt.Printf("Best match found: UserID=%d\n", match.UserID)
	} else {
		fmt.Println("No match found")
	}

	// Delete a user
	search.Delete(1)

	// Clear all data
	search.Close()
}
```

---

## API

### `func Create(UserID uint64, Language string, YourStart, YourEnd, YourSex, MyAge, MySex int, Priority bool, interests ...string)`

Creates and stores a new user profile.

### `func Search(MyID uint64, Language string, YourStart, YourEnd, YourSex, MyAge, MySex int, interests ...string) *SearchEntry`

Searches for the best matching user.
Returns `nil` if no suitable match is found.

### `func Delete(UserID uint64)`

Deletes a user from the store.

### `func Close()`

Clears the in-memory store.

---

## Testing

Run all tests:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

Example results (Apple M4, Go 1.22):

```
BenchmarkCreate-10                          	 4105684	       333.8 ns/op	     307 B/op	       2 allocs/op
BenchmarkSearchSingleUser-10                	132750019	         9.105 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchMultipleUsers-10             	 1000000	      1091 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchWithManyInterests-10         	  964914	      1101 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchWithDifferentLanguages-10    	 1423964	       868.3 ns/op	       0 B/op	       0 allocs/op
```

---