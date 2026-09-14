package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Panicln(err)
	}

	if code > 499 {
		log.Printf("Responding with 5xx error: %s", err)
	}
	type jsonErr struct {
		Err string `json:"error"`
	}
	respondWithJson(w, code, jsonErr{
		Err: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Somethung went wrong: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}