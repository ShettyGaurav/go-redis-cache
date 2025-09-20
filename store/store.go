package store

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type ValueType int

const (
	StringType ValueType = iota
	ListType
	SetType
	HashType
)

type RedisValue struct {
	Type       ValueType
	StringVal  string
	ListVal    []string
	SetVal     map[string]struct{}
	HashVal    map[string]string
	Expiration int64
}
type RedisStore struct {
	data     map[string]*RedisValue
	mu       sync.RWMutex
	filename string
}

func (rs *RedisStore) startGC() {
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for {
			<-ticker.C
			rs.mu.Lock()
			now := time.Now().Unix()
			for k, v := range rs.data {
				if v.Expiration > 0 && now > v.Expiration {
					delete(rs.data, k)
				}
			}
			rs.mu.Unlock()
		}
	}()
}

func NewRedisStore(filename string) *RedisStore {
	rs := &RedisStore{
		data:     make(map[string]*RedisValue),
		filename: filename,
	}
	rs.loadFromFile()
	rs.startGC()
	return rs
}

func (rs *RedisStore) Set(key string, value string, expiration time.Duration) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).Unix()
	}
	rs.data[key] = &RedisValue{
		Type:       StringType,
		StringVal:  value,
		Expiration: exp,
	}
	rs.appendToFile(key, value, exp, "SET")
}

func (rs *RedisStore) Get(key string) (string, bool) {
	rs.mu.RLock()
	item, found := rs.data[key]
	rs.mu.RUnlock()
	if !found {
		return "", false
	}
	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		rs.mu.Lock()
		if storedItem, exists := rs.data[key]; exists && storedItem == item {
			delete(rs.data, key)
		}
		rs.mu.Unlock()
		return "", false
	}
	return item.StringVal, true
}

func (rs *RedisStore) Delete(key string) int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if _, found := rs.data[key]; found {
		delete(rs.data, key)
		return 1
	}
	return 0
}

func (rs *RedisStore) Exists(key string) bool {
	rs.mu.RLock()
	item, found := rs.data[key]
	rs.mu.RUnlock()
	if !found {
		return false
	}

	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		rs.mu.Lock()
		defer rs.mu.Unlock()
		if storedItem, exists := rs.data[key]; exists && storedItem == item {
			delete(rs.data, key)
			return false
		}
		return false
	}
	return true
}

func (rs *RedisStore) loadFromFile() {
	f, err := os.Open(rs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Error loading from file:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var cmd, key, value string
		var expiration int64
		line := scanner.Text()

		_, err := fmt.Sscanf(line, "%s %s %s %d", &cmd, &key, &value, &expiration)
		if err != nil {
			fmt.Println("Failed to parse line:", line, err)
			continue
		}

		if expiration > 0 && time.Now().Unix() > expiration {
			continue
		}

		switch cmd {
		case "SET":
			rs.data[key] = &RedisValue{
				Type:       StringType,
				StringVal:  value,
				Expiration: expiration,
			}
		case "LPUSH":
			item, exists := rs.data[key]
			if !exists {
				item = &RedisValue{
					Type:    ListType,
					ListVal: []string{},
				}
				rs.data[key] = item
			}
			if item.Type == ListType {
				item.ListVal = append([]string{value}, item.ListVal...)
				item.Expiration = expiration
			}
		case "LPOP":
			if item, exists := rs.data[key]; exists && item.Type == ListType && len(item.ListVal) > 0 {
				item.ListVal = item.ListVal[1:]
				if len(item.ListVal) == 0 {
					delete(rs.data, key)
				}
				item.Expiration = expiration
			}
		}
	}
}

func (rs *RedisStore) appendToFile(key, value string, expiration int64, cmd string) {
	f, err := os.OpenFile(rs.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error appending to file:", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s %s %s %d\n", cmd, key, value, expiration))
	if err != nil {
		fmt.Println("Error appending to file:", err)
	}
}

func (rs *RedisStore) GetType(key string) (ValueType, bool) {
	rs.mu.RLock()
	item, found := rs.data[key]
	rs.mu.RUnlock()

	if !found {
		return StringType, false
	}

	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		rs.mu.Lock()
		defer rs.mu.Unlock()
		if storedItem, exists := rs.data[key]; exists && storedItem == item {
			delete(rs.data, key)
		}
		return StringType, false
	}

	return item.Type, true
}

func (rs *RedisStore) LPush(key string, value string) int {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	item, exists := rs.data[key]
	if !exists {
		item = &RedisValue{
			Type:    ListType,
			ListVal: []string{},
		}
		rs.data[key] = item
	} else if item.Type != ListType {
		return 0
	}
	item.ListVal = append([]string{value}, item.ListVal...)
	rs.appendToFile(key, value, item.Expiration, "LPUSH")
	return len(item.ListVal)
}

func (rs *RedisStore) LPop(key string) (string, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	item, exists := rs.data[key]
	if !exists {
		return "", false
	}
	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		delete(rs.data, key)
		return "", false
	}
	if item.Type != ListType {
		return "", false
	}
	if len(item.ListVal) == 0 {
		return "", false
	}
	value := item.ListVal[0]
	item.ListVal = item.ListVal[1:]
	if len(item.ListVal) == 0 {
		delete(rs.data, key)
	}
	rs.appendToFile(key, value, item.Expiration, "LPOP")
	return value, true
}

func (rs *RedisStore) LRange(key string, start, stop int) ([]string, bool) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	item, exists := rs.data[key]
	if !exists {
		return nil, false
	}
	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		return nil, false
	}
	if item.Type != ListType {
		return nil, false
	}

	list := item.ListVal
	length := len(list)

	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		return []string{}, true
	}

	return list[start : stop+1], true
}
