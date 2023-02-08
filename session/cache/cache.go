package cache

import "time"

func Set(key string, value interface{}, lifeTime time.Duration) (err error) {
	cache := si.GetCache()
	err = cache.Add(key, value, lifeTime)
	return
}

func Get(key string) (data interface{}, found bool) {
	cache := si.GetCache()
	data, found = cache.Get(key)
	return
}

func Del(key string) {
	cache := si.GetCache()
	cache.Delete(key)
}
