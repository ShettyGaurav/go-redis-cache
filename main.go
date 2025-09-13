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
	store.Set("Goo", "too", 20*time.Hour)

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

	//-----------------LIST Operations-----------------
	store.LPush("mylist", "value1")
	store.LPush("mylist", "value2")
	store.LPush("mylist", "value3")
	store.LPush("mylist", "value4")
	store.LPush("mylist", "value5")

	store.LRange("mylist", 0, 2)
	if values, ok := store.LRange("mylist", 0, 4); ok {
		fmt.Println("Values:", values)
	} else {
		fmt.Println("Values not found")
	}

	store.LPop("mylist")
	if value, ok := store.LPop("mylist"); ok {
		fmt.Println("Popped value:", value)
	} else {
		fmt.Println("Popped value not found")
	}

}
