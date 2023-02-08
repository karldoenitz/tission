package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var cacheManager *cache.Cache

func produceCacheManager(defaultExpiration, cleanupInterval time.Duration) *cache.Cache {
	return cache.New(defaultExpiration, cleanupInterval)
}
