package cache

import (
	"cache-demo/internal/types"
	"log"
	"sync"
	"time"
)

type Cache struct {
	Data           map[types.GridKey]types.ForecastVals
	mu             *sync.RWMutex
	BucketDeadline time.Time
}

func CreateCache() *Cache {
	cache := make(map[types.GridKey]types.ForecastVals)
	var mu sync.RWMutex
	deadline := time.Now().Add(-time.Hour)
	return &Cache{cache, &mu, deadline}
}

func (c *Cache) Get(instance types.GridKey) (types.ForecastVals, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.Data[instance]
	if exists {
		log.Println("Cache Used")
	}
	return val, exists
}

func (c *Cache) Set(instance types.GridKey, value types.ForecastVals) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Data[instance] = value
	log.Println("Cache Instance Created")
}
