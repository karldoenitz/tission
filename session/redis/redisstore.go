package redis

import (
	"fmt"
	"time"

	"github.com/alexedwards/scs/stores/redisstore"
	"github.com/gomodule/redigo/redis"
)

// Store represents the currently configured session session store. It is essentially
// a wrapper around a redisstore's RedisStore
type Store struct {
	*redisstore.RedisStore
	pool *redis.Pool
}

// New returns a new RedisStore instance. The pool parameter should be a pointer to a
// Redigo connection pool. See https://godoc.org/github.com/garyburd/redigo/redis#Pool.
func NewStore(pool *redis.Pool) *Store {
	return &Store{redisstore.New(pool), pool}
}

// Save adds a session token and data to the RedisStore instance with the given expiry time.
// If the session token already exists then the data and expiry time are updated.
func (r *Store) Save(token string, b []byte, expiry time.Time) error {
	conn := r.pool.Get()
	defer conn.Close()
	fmt.Printf("%s lide seconds is :%d\n", redisstore.Prefix+token, int(expiry.Sub(time.Now()).Seconds()))
	_, err := conn.Do("SETEX", redisstore.Prefix+token, int(expiry.Sub(time.Now()).Seconds()), b)
	return err
}
