package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"github.com/montruh-afk/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	Db             *database.Queries
	Platform       string
}

const (
	metrics = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) RequestHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, metrics, cfg.FileServerHits.Load()); err != nil {
		log.Printf("Something broke on out end: %s", err)
	}
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You do not have the permission to interact via this endpoint"))
		return
	}
	ctx := r.Context()
	if err := cfg.Db.ResetUsers(ctx); err != nil {
		log.Panicf("Something went wrong: %s", err)
	}
	cfg.FileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type getuser struct {
		Email string `json:"email"`
	}

	createUser := getuser{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&createUser); err != nil {
		log.Printf("Something broke while attempting to decode json: %s", err)
		return
	}
	if !strings.Contains(createUser.Email, "@") && !strings.Contains(strings.ToLower(createUser.Email), ".co") {
		respondWithError(w, 500, "Invalid email", nil)
		return
	}
	user, err := cfg.Db.CreateUser(ctx, createUser.Email)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		respondWithError(w, 500, "Something broke while reaching our records", err)
		return
	}

	respondWithJson(w, http.StatusCreated, user)

}

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	params, err := validateChirp(r)
	if err != nil {
		log.Fatal(err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	dbChirp := database.CreateChirpParams{
		UserID: params.UserID,
		Body: params.Body,
	}
	
	ctx := r.Context()
	data, err := cfg.Db.CreateChirp(ctx, dbChirp)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}
	
	respondWithJson(w, http.StatusCreated, data)
}
