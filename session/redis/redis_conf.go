package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var redisPool *redis.Pool

func produceRedisPool(addr string, maxIdle, timeout int, auth string, dbNo interface{}, pwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr, redis.DialPassword(pwd), redis.DialDatabase(dbNo.(int)))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
