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
	test "myredis/tests"
	"testing"
	"time"
)

func main() {
	rs := store.NewRedisStore("data.txt")

	// -------- Test loadFromFile --------
	if val, ok := rs.Get("foo"); ok {
		fmt.Println("Reloaded foo from file:", val)
	} else {
		fmt.Println("foo not found in file")
	}

	// -------- Test LPUSH --------
	rs.LPush("mylist", "value1")
	rs.LPush("mylist", "value2")
	rs.LPush("mylist", "value3")

	if values, ok := rs.LRange("mylist", 0, 10); ok {
		fmt.Println("List after LPush:", values)
	}

	// -------- Test LPOP --------
	if val, ok := rs.LPop("mylist"); ok {
		fmt.Println("Popped from mylist:", val)
	}

	if values, ok := rs.LRange("mylist", 0, 10); ok {
		fmt.Println("List after LPop:", values)
	}

	// -------- Test Persistence --------
	fmt.Println("Restart program to see if foo/mylist reload from file.")
	time.Sleep(2 * time.Second)
	rs.Set("foo", "bar", 20*time.Hour)

	test.TestEncodeBulkString(&testing.T{})
	test.TestEncodeError(&testing.T{})
	test.TestEncodeInteger(&testing.T{})
	test.TestEncodeSimpleString(&testing.T{})
}
