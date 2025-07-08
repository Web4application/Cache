package core

import (
    "encoding/json"
    "errors"
    "os"
    "sync"
    "time"
)

// Item represents a cached value with optional expiry.
type Item struct {
    Value  string    `json:"value"`
    Expiry time.Time `json:"expiry,omitempty"`
}

// Engine is a simple in‑memory cache with TTL and JSON file backup.
type Engine struct {
    mu    sync.RWMutex
    store map[string]Item
    // persistent file path
    filePath string
}

// NewEngine creates a new cache engine and optionally loads from file.
func NewEngine(filePath string) *Engine {
    e := &Engine{
        store:    make(map[string]Item),
        filePath: filePath,
    }
    if filePath != "" {
        _ = e.LoadFromFile()
    }
    return e
}

// Set inserts/updates a key with optional TTL in seconds (0 = no expiry).
func (e *Engine) Set(key, value string, ttlSeconds int) {
    e.mu.Lock()
    defer e.mu.Unlock()
    var expiry time.Time
    if ttlSeconds > 0 {
        expiry = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
    }
    e.store[key] = Item{Value: value, Expiry: expiry}
}

// Get retrieves a key; returns value, found flag.
func (e *Engine) Get(key string) (string, bool) {
    e.mu.RLock()
    item, ok := e.store[key]
    e.mu.RUnlock()
    if !ok {
        return "", false
    }
    if !item.Expiry.IsZero() && time.Now().After(item.Expiry) {
        // expired
        e.Delete(key)
        return "", false
    }
    return item.Value, true
}

// Delete removes a key.
func (e *Engine) Delete(key string) {
    e.mu.Lock()
    delete(e.store, key)
    e.mu.Unlock()
}

// List returns all non‑expired keys (debug use).
func (e *Engine) List() map[string]Item {
    e.mu.RLock()
    defer e.mu.RUnlock()
    result := make(map[string]Item)
    now := time.Now()
    for k, v := range e.store {
        if v.Expiry.IsZero() || now.Before(v.Expiry) {
            result[k] = v
        }
    }
    return result
}

// SaveToFile writes cache to disk.
func (e *Engine) SaveToFile() error {
    if e.filePath == "" {
        return errors.New("no file path configured")
    }
    e.mu.RLock()
    data, err := json.MarshalIndent(e.store, "", "  ")
    e.mu.RUnlock()
    if err != nil {
        return err
    }
    if err := os.MkdirAll(os.DirFS(e.filePath).(string), 0o755); err != nil {
        // ensure dir exists
    }
    return os.WriteFile(e.filePath, data, 0o644)
}

// LoadFromFile loads cache from disk (ignores expired entries).
func (e *Engine) LoadFromFile() error {
    if e.filePath == "" {
        return errors.New("no file path configured")
    }
    data, err := os.ReadFile(e.filePath)
    if err != nil {
        return err
    }
    var m map[string]Item
    if err := json.Unmarshal(data, &m); err != nil {
        return err
    }
    now := time.Now()
    e.mu.Lock()
    for k, v := range m {
        if v.Expiry.IsZero() || now.Before(v.Expiry) {
            e.store[k] = v
        }
    }
    e.mu.Unlock()
    return nil
}
