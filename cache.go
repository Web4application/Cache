package main

import (
    "sync"
    "time"
)

type CacheItem struct {
    Value      string
    ExpiryTime time.Time
}

type Cache struct {
    items map[string]CacheItem
    mu    sync.RWMutex
    ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
    c := &Cache{
        items: make(map[string]CacheItem),
        ttl:   ttl,
    }
    go c.evictExpired()
    return c
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = CacheItem{
        Value:      value,
        ExpiryTime: time.Now().Add(c.ttl),
    }
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, exists := c.items[key]
    if !exists || time.Now().After(item.ExpiryTime) {
        return "", false
    }
    return item.Value, true
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

func (c *Cache) evictExpired() {
    for {
        time.Sleep(c.ttl)
        c.mu.Lock()
        for k, v := range c.items {
            if time.Now().After(v.ExpiryTime) {
                delete(c.items, k)
            }
        }
        c.mu.Unlock()
    }
}
