package main

import (
	"net/http"
	"github.com/montruh-afk/chirpy/internal"
)

func startUp(handler *http.ServeMux, cfg *internal.ApiConfig) {
	handler.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./web")))))
	handler.HandleFunc("GET /admin/metrics", cfg.RequestHits)
	handler.HandleFunc("GET /api/healthz", internal.ReadinessEndpoint)
	handler.HandleFunc("POST /admin/reset", cfg.Reset)
	handler.HandleFunc("POST /api/users", cfg.CreateUser)
	handler.HandleFunc("POST /api/chirps", cfg.CreateChirp)
	handler.HandleFunc("GET /api/chirps", cfg.GetChirps)
	handler.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirp)
	handler.HandleFunc("POST /api/login", cfg.HandlerLogin)
}