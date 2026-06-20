package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/coocood/freecache"
)

type InMemoryCache struct {
	cache *freecache.Cache
}

func NewInMemory() *InMemoryCache {
	return &InMemoryCache{
		cache: freecache.NewCache(100 * 1024 * 1024),
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.cache.Get([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("cache miss: %w", err)
	}
	return val, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	expire := int(ttl.Seconds())
	if expire <= 0 {
		expire = 300
	}
	return c.cache.Set([]byte(key), value, expire)
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.cache.Del([]byte(key))
	return nil
}

func (c *InMemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.cache.Get([]byte(key))
	return err == nil, nil
}

func (c *InMemoryCache) Flush(ctx context.Context) error {
	c.cache.Clear()
	return nil
}

func (c *InMemoryCache) Close() error {
	c.cache.Clear()
	return nil
}
