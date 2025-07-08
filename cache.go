package main

import (
    "container/list"
    "sync"
    "time"
)

type CacheItem struct {
    Key        string
    Value      string
    ExpiryTime time.Time
}

type Cache struct {
    items     map[string]*list.Element
    order     *list.List
    mu        sync.Mutex
    ttl       time.Duration
    capacity  int
    persist   bool
}

func NewCache(ttl time.Duration, capacity int) *Cache {
    c := &Cache{
        items:    make(map[string]*list.Element),
        order:    list.New(),
        ttl:      ttl,
        capacity: capacity,
    }
    go c.cleanup()
    return c
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if el, ok := c.items[key]; ok {
        c.order.MoveToFront(el)
        el.Value.(*CacheItem).Value = value
        el.Value.(*CacheItem).ExpiryTime = time.Now().Add(c.ttl)
        return
    }

    if c.order.Len() >= c.capacity {
        oldest := c.order.Back()
        if oldest != nil {
            c.order.Remove(oldest)
            delete(c.items, oldest.Value.(*CacheItem).Key)
        }
    }

    item := &CacheItem{Key: key, Value: value, ExpiryTime: time.Now().Add(c.ttl)}
    el := c.order.PushFront(item)
    c.items[key] = el
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    el, ok := c.items[key]
    if !ok || time.Now().After(el.Value.(*CacheItem).ExpiryTime) {
        return "", false
    }
    c.order.MoveToFront(el)
    return el.Value.(*CacheItem).Value, true
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if el, ok := c.items[key]; ok {
        c.order.Remove(el)
        delete(c.items, key)
    }
}

func (c *Cache) cleanup() {
    for {
        time.Sleep(c.ttl)
        c.mu.Lock()
        for k, el := range c.items {
            if time.Now().After(el.Value.(*CacheItem).ExpiryTime) {
                c.order.Remove(el)
                delete(c.items, k)
            }
        }
        c.mu.Unlock()
    }
}
