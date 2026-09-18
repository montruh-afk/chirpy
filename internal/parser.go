package internal

import (
	"strings"
	"time"
	"net/http"
)



func checkProfane(body string) string {
	forbidden := map[string]any{
		"kerfuffle": nil, 
		"sharbert": nil, 
		"fornax": nil,
	}

	cleaned_body := ""

	slice := strings.SplitSeq(body, " ")
	for word := range slice{
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