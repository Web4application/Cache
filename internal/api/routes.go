package api

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/Web4application/cache/internal/core"
)

type API struct {
    Engine *core.Engine
}

func (a *API) handleSet(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    value := r.URL.Query().Get("value")
    ttlStr := r.URL.Query().Get("ttl")
    if key == "" {
        http.Error(w, "key required", http.StatusBadRequest)
        return
    }
    ttl := 0
    if ttlStr != "" {
        if n, err := strconv.Atoi(ttlStr); err == nil && n > 0 {
            ttl = n
        }
    }
    a.Engine.Set(key, value, ttl)
    w.WriteHeader(http.StatusCreated)
}

func (a *API) handleGet(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key == "" {
        http.Error(w, "key required", http.StatusBadRequest)
        return
    }
    val, ok := a.Engine.Get(key)
    if !ok {
        http.NotFound(w, r)
        return
    }
    w.Write([]byte(val))
}

func (a *API) handleDelete(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key == "" {
        http.Error(w, "key required", http.StatusBadRequest)
        return
    }
    a.Engine.Delete(key)
    w.WriteHeader(http.StatusNoContent)
}

func (a *API) handleBackup(w http.ResponseWriter, _ *http.Request) {
    if err := a.Engine.SaveToFile(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (a *API) handleList(w http.ResponseWriter, _ *http.Request) {
    m := a.Engine.List()
    _ = json.NewEncoder(w).Encode(m)
}

func (a *API) Register(mux *http.ServeMux) {
    mux.HandleFunc("/set", a.handleSet)
    mux.HandleFunc("/get", a.handleGet)
    mux.HandleFunc("/delete", a.handleDelete)
    mux.HandleFunc("/backup", a.handleBackup)
    mux.HandleFunc("/list", a.handleList)
}
