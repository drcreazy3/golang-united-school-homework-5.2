package cache

import (
	"time"
)

const CLEAR_INTERVAL = 1

type Cache struct {
	elementMap             map[string]CacheElement
	clearExpiredElementsAt time.Time
}

type CacheElement struct {
	value    string
	deadline time.Time
}

func NewCache() Cache {
	c := Cache{}
	c.elementMap = map[string]CacheElement{}
	c.clearExpiredElementsAt = time.Now().Add(time.Second * time.Duration(CLEAR_INTERVAL))
	return c
}

func (c Cache) Get(key string) (string, bool) {
	c.checkDeadline()
	v := ""
	ok := false
	if cacheElement, ok := c.elementMap[key]; ok {
		if cacheElement.deadline.IsZero() || time.Now().Before(cacheElement.deadline) {
			return cacheElement.value, ok
		}
	}

	return v, ok
}

func (c Cache) Put(key, value string) {
	c.checkDeadline()
	c.elementMap[key] = CacheElement{value: value}
}

func (c Cache) Keys() []string {
	c.checkDeadline()
	var keys []string
	for k := range c.elementMap {
		keys = append(keys, k)
	}

	return keys
}

func (c Cache) PutTill(key, value string, deadline time.Time) {
	if deadline.Before(c.clearExpiredElementsAt) {
		c.clearExpiredElementsAt = deadline
	}
	c.elementMap[key] = CacheElement{value: value, deadline: deadline}
	c.checkDeadline()
}

func (c Cache) checkDeadline() {
	if time.Now().After(c.clearExpiredElementsAt) {
		c.clearExpiredElementsAt = time.Now().Add(time.Second * time.Duration(CLEAR_INTERVAL))
		for k, cacheElement := range c.elementMap {
			if cacheElement.deadline.IsZero() {
				continue
			}

			if time.Now().After(cacheElement.deadline) {
				delete(c.elementMap, k)
				continue
			}

			if cacheElement.deadline.Before(c.clearExpiredElementsAt) {
				c.clearExpiredElementsAt = cacheElement.deadline
			}
		}
	}
}
