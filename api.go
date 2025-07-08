package main

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
)

type KV struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

func setupRoutes(c *Cache) *mux.Router {
    r := mux.NewRouter()

    r.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
        var kv KV
        json.NewDecoder(r.Body).Decode(&kv)
        c.Set(kv.Key, kv.Value)
        w.Write([]byte("OK"))
    }).Methods("POST")

    r.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        value, found := c.Get(key)
        if !found {
            http.Error(w, "Not found", http.StatusNotFound)
            return
        }
        json.NewEncoder(w).Encode(KV{Key: key, Value: value})
    }).Methods("GET")

    r.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Query().Get("key")
        c.Delete(key)
        w.Write([]byte("Deleted"))
    }).Methods("DELETE")

    return r
}
