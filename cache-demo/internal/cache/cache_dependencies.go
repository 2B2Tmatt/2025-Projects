package cache

import (
	"cache-demo/internal/types"
	"sync"
)

type Cache struct {
	Data map[types.GridKey]types.ForecastVals
	mu   *sync.RWMutex
}

func CreateCache() *Cache {
	cache := make(map[types.GridKey]types.ForecastVals)
	var mu sync.RWMutex
	return &Cache{cache, &mu}
}

func (c *Cache) Get(instance types.GridKey) (types.ForecastVals, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.Data[instance]
	return val, exists
}

func (c *Cache) Set(instance types.GridKey, value types.ForecastVals) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Data[instance] = value
}

func (c *Cache) Delete(instance types.GridKey) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Data, instance)
}
