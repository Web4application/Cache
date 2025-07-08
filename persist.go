package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
)

func SaveToFile(cache *Cache, filename string) error {
    cache.mu.Lock()
    defer cache.mu.Unlock()

    data := make(map[string]CacheItem)
    for k, el := range cache.items {
        data[k] = *el.Value.(*CacheItem)
    }

    jsonBytes, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    return ioutil.WriteFile(filename, jsonBytes, 0644)
}

func LoadFromFile(cache *Cache, filename string) error {
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return nil // File doesn't exist yet
    }

    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }

    items := make(map[string]CacheItem)
    if err := json.Unmarshal(data, &items); err != nil {
        return err
    }

    cache.mu.Lock()
    defer cache.mu.Unlock()

    for _, item := range items {
        el := cache.order.PushFront(&item)
        cache.items[item.Key] = el
    }

    return nil
}
