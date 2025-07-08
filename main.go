package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "gopkg.in/yaml.v3"

    "github.com/Web4application/cache/internal/api"
    "github.com/Web4application/cache/internal/core"
)

type Config struct {
    Server struct {
        Port int `yaml:"port"`
    } `yaml:"server"`
    Storage struct {
        BackupFile        string `yaml:"backupFile"`
        AutoBackupSeconds int    `yaml:"autoBackupSeconds"`
    } `yaml:"storage"`
}

func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}

func main() {
    cfg, err := loadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("config load: %v", err)
    }
    eng := core.NewEngine(cfg.Storage.BackupFile)
    apiSrv := &api.API{Engine: eng}

    mux := http.NewServeMux()
    apiSrv.Register(mux)
    mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // auto backup ticker
    if cfg.Storage.AutoBackupSeconds > 0 {
        ticker := time.NewTicker(time.Duration(cfg.Storage.AutoBackupSeconds) * time.Second)
        go func() {
            for range ticker.C {
                if err := eng.SaveToFile(); err != nil {
                    log.Printf("backup error: %v", err)
                }
            }
        }()
    }

    addr := ":" + strconv.Itoa(cfg.Server.Port)
    log.Printf("Cache service listening at %s", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}
