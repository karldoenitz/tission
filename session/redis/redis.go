package redis

import (
	"encoding/json"
	"fmt"
	"time"
)

const DbCacheTime = 3600 * time.Second

func Set(key string, value interface{}, lifeTime time.Duration) (error, []byte) {
	redisPool := si.GetRedisPool()
	conn := redisPool.Get()
	defer conn.Close()

	var (
		err error
		v   []byte
	)
	switch value.(type) {
	case string:
		v = []byte(value.(string))
	case []byte:
		v = value.([]byte)
	default:
		v, err = json.Marshal(value)
	}
	s := lifeTime.Seconds()
	_, err = conn.Do("SETEX", key, int(s), v)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("SETEX key[%s] failed", key)
	}
	return err, v
}

func Get(key string) (data []byte, found bool) {
	redisPool := si.GetRedisPool()
	conn := redisPool.Get()
	defer conn.Close()
	reply, _ := conn.Do("GET", key)
	if reply != nil {
		data = reply.([]byte)
		found = true
	} else {
		found = false
	}
	return
}

func Del(key string) int64 {
	redisPool := si.GetRedisPool()
	conn := redisPool.Get()
	defer conn.Close()
	reply, _ := conn.Do("DEL", key)
	return reply.(int64)
}
