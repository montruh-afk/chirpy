package internal

import (
	"fmt"
	"net/http"
	"sync/atomic"
)
type ApiConfig struct {
	FileServerHits atomic.Int32
	file32 int32
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cfg.file32 = cfg.FileServerHits.Add(1)
	next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) RequestHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.file32)))
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits = atomic.Int32{}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}