package main

import (
	"github.com/montruh-afk/chirpy/internal"
	"net/http"
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
	handler.HandleFunc("POST /api/revoke", cfg.Revoke)
	handler.HandleFunc("POST /api/refresh", cfg.Refresh)
	handler.HandleFunc("PUT /api/users", cfg.UpdateUserLogin)
}
