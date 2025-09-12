// package main

// import (
// 	"fmt"
// 	"myredis/store"
// 	"time"
// )

// func main() {
// 	store := store.NewRedisStore("data.txt")
// 	val1, ok1 := store.Get("foo")
// 	if ok1 {
// 		fmt.Println("Reloaded foo:", val1)
// 	} else {
// 		fmt.Println("foo not found")
// 	}
// 	store.Set("foo", "bar", 20*time.Hour)
// 	val, ok := store.Get("foo")
// 	if ok {
// 		fmt.Println("Value:", val)
// 	} else {
// 		fmt.Println("Value not found")
// 	}
// 	store.Set("goo", "too", 10*time.Second)
// 	fmt.Println("Value exists:", store.Exists("foo"))
// }

package main

import (
	"fmt"
	"myredis/store"
	"time"
)

func main() {
	store := store.NewRedisStore("data.txt")

	val1, ok1 := store.Get("foo")
	if ok1 {
		fmt.Println("Reloaded foo:", val1)
	} else {
		fmt.Println("foo not found")
	}

	store.Set("foo", "bar", 20*time.Hour)

	val, ok := store.Get("foo")
	if ok {
		fmt.Println("Value:", val)
	} else {
		fmt.Println("Value not found")
	}

	store.Set("goo", "too", 10*time.Second)

	exists := store.Exists("foo")
	fmt.Println("Value exists:", exists)

	deleted := store.Delete("goo")
	fmt.Println("Deleted goo:", deleted)

	typ, found := store.GetType("foo")
	if found {
		fmt.Println("Type of foo:", typ)
	} else {
		fmt.Println("foo not found or expired")
	}
}
