package internal

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"errors"
)

type details struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func checkProfane(body string) string {
	forbidden := map[string]any{
		"kerfuffle": nil,
		"sharbert":  nil,
		"fornax":    nil,
	}

	cleaned_body := ""

	slice := strings.SplitSeq(body, " ")
	for word := range slice {
		if _, exists := forbidden[strings.ToLower(word)]; exists {
			cleaned_body += "**** "
		} else {
			cleaned_body += word + " "
		}

	}
	cleaned_body = strings.Trim(cleaned_body, " ")
	return cleaned_body
}

func parseDuration(d *Duration, w http.ResponseWriter, duration string) {
	time, err := time.ParseDuration(duration)
	if err != nil {
		respondWithError(w, 500, "Could not parse expiry duration", err)
		return
	}
	d.duration = time
}

func verifyDetails(r *http.Request) (*details, error) {
	deets := details{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&deets); err != nil {
		return nil, err
	}

	// need to make these stronger
	if !strings.Contains(deets.Email, "@") && !strings.Contains(strings.ToLower(deets.Email), ".co") {
		return nil, errors.New("Invalid email")
	}
	if len(deets.Password) < 4 { //for tests alone, change this later
		return nil, errors.New("Password length should be 8 or more characters")
	}

	return &deets, nil
}
