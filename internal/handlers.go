package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/google/uuid"
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
		log.Println(err)
		return chirp{}, fmt.Errorf("Something went wrong while attempting to decode json: %s", err)
	}

	if len(params.Body) > maxChirpLength {
		return params, fmt.Errorf("Chirp is too long\n")
	}

	if len(params.Body) <= maxChirpLength {
		params.Body = checkProfane(params.Body)
	}
	return params, nil
}
