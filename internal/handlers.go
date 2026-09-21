package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"github.com/montruh-afk/chirpy/internal/database"
	"context"
)

const (
	maxChirpLength = 140
)

type chirp struct {
	Body   string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func ReadinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateChirp(r *http.Request) (chirp, error) {
	params := chirp{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		return params, fmt.Errorf("Something went wrong while attempting to decode json: %s", err)
	}

	if len(params.Body) > maxChirpLength {
		return params, fmt.Errorf("Chirp is too long\n")
	}

	if len(params.Body) <= maxChirpLength {
		params.Body = checkProfane(params.Body)
	}
	return params, nil
}

func fetchChirp(cfg *ApiConfig, r *http.Request) (database.Chirp, error) {
	chirpID := r.PathValue("chirpID")
	if len(chirpID) < 1 {
		return  database.Chirp{}, fmt.Errorf("Required parameter omitted")
	}
	if err := uuid.Validate(chirpID); err != nil {
		return database.Chirp{}, fmt.Errorf("Invalid id: %s", err)
	}
	id, err := uuid.Parse(chirpID)
	if err != nil {
		return database.Chirp{}, fmt.Errorf("Could not verify the id provided: %v", err)
	}
	chirp, err := cfg.Db.GetChirp(r.Context(), id)
	if err != nil {
		return database.Chirp{}, fmt.Errorf("Chirp with ID %v not found: %v", chirpID, err)
	}
	return chirp, nil
}

func fetchUserChirps(cfg *ApiConfig, author string, ctx context.Context) ([]database.Chirp, error) {
	if err := uuid.Validate(author); err != nil {
		return []database.Chirp{}, err
	}
	id, err := uuid.Parse(author)
	if err != nil {
		return []database.Chirp{}, err
	}
	chirps, err := cfg.Db.GetChirpsByUser(ctx, id)
	return chirps, nil
}