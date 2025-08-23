package cache

import (
	"cache-demo/internal/types"
	"log"
	"time"
)

func (c *Cache) BackgroundProcessor() {
	for {
		c.mu.RLock()
		validDeadline := c.BucketDeadline.After(time.Now())
		c.mu.RUnlock()
		if validDeadline {
			c.SleepUntilDeadline()
			c.Flush()
			c.mu.Lock()
			c.BucketDeadline = time.Now().Add(-time.Hour)
			c.mu.Unlock()
			log.Println("Cache flushed")
		} else {
			c.mu.Lock()
			for value := range c.Data {
				c.BucketDeadline = c.Data[value].ExpiresAt
				break
			}
			c.mu.Unlock()
			time.Sleep(time.Minute)
		}
	}
}

func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Data = make(map[types.GridKey]types.ForecastVals)
}

func (c *Cache) SleepUntilDeadline() {
	c.mu.RLock()
	sleepDuration := time.Until(c.BucketDeadline)
	c.mu.RUnlock()
	time.Sleep(sleepDuration)
}
