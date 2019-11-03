package redis

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"runtime"
	"time"
)

const DbCacheTime = 3600 * time.Second

func SaveString(key string, value string, lifeTime time.Duration) error {
	redisPool := si.GetRedisPool()
	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, int(lifeTime.Seconds()), []byte(value))

	return err
}

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

func Hset(key, field, value string) (err error) {
	redisPool := si.GetRedisPool()
	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("HSET", key, field, value)
	return err
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

func GetInterface(key string) (reply interface{}, found bool) {
	redisPool := si.GetRedisPool()
	conn := redisPool.Get()
	defer conn.Close()
	reply, _ = conn.Do("GET", key)
	if reply != nil {
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

func Exists(key string) (value bool, err error) {
	redisPool := si.GetRedisPool()
	conn := redisPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

func makeMillisecondTimestamp(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func GenKey(params ...interface{}) (key string, err error) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("gen redis key error: caller error")
		return
	}
	name := runtime.FuncForPC(pc).Name()
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Printf("json marshal runtime-pc-name error")
		return
	}
	key = fmt.Sprintf("%x", md5.Sum(append(b, []byte(name)...)))
	return
}

func GetOrSet(key string, f func() (interface{}, time.Duration, error)) (data []byte, err error) {
	data, found := Get(key)
	if found {
		return data, nil
	}
	setData, liftTime, err := f()
	if err == nil {
		err, data = Set(key, setData, liftTime)
	}
	return data, err
}
