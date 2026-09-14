package internal

import (
	"encoding/json"
	"net/http"
)

const (
	maxChirpLength = 140
)

func ReadinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	// temporary struct to decode request body (alternative to casting []byte to string then extracting string value after the "value" title )
	type params struct {
		Body string `json:"body"`
	}

	type isValid struct {
		Valid bool `json:"valid"`
	}

	parameters := params{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&parameters); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while attempting to decode json", err)
	}
	if len(parameters.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return

	} else if len(parameters.Body) <= maxChirpLength{
		checkProfane(w, parameters.Body)
		return
	}
}
