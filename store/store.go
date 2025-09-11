package store

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type RedisValue struct {
	Data       string
	Expiration int64
}
type redisStore struct {
	data     map[string]*RedisValue
	mu       sync.RWMutex
	filename string
}

func (rs *redisStore) startGC() {
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

func NewRedisStore(filename string) *redisStore {
	rs := &redisStore{
		data:     make(map[string]*RedisValue),
		filename: filename,
	}
	rs.loadFromFile()
	rs.startGC()
	return rs
}

func (rs *redisStore) Set(key string, value string, expiration time.Duration) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).Unix()
	}
	rs.data[key] = &RedisValue{
		Data:       value,
		Expiration: exp,
	}
	rs.appendToFile(key, value, exp, "SET")
}

func (rs *redisStore) Get(key string) (string, bool) {
	rs.mu.RLock()
	item, found := rs.data[key]
	rs.mu.RUnlock()
	if !found {
		return "", false
	}
	if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
		rs.mu.Lock()
		delete(rs.data, key)
		rs.mu.Unlock()
		return "", false
	}
	return item.Data, true
}

func (rs *redisStore) Delete(key string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	delete(rs.data, key)
}

func (rs *redisStore) Exists(key string) bool {
	rs.mu.RLock()
	_, found := rs.data[key]
	if found {
		item, found := rs.data[key]
		if found {
			if item.Expiration > 0 && time.Now().Unix() > item.Expiration {
				delete(rs.data, key)
				found = false
			}
		}
	}
	rs.mu.RUnlock()
	return found
}

func (rs *redisStore) loadFromFile() {
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

		// skip expired keys
		if expiration == 0 || time.Now().Unix() <= expiration {
			rs.data[key] = &RedisValue{
				Data:       value,
				Expiration: expiration,
			}
		}
	}
}

func (rs *redisStore) appendToFile(key, value string, expiration int64, cmd string) {
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
