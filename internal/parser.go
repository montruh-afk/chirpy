package internal

import (
	"net/http"
	"strings"
)



func checkProfane(w http.ResponseWriter, body string) {
	forbidden := map[string]any{
		"kerfuffle": nil, 
		"sharbert": nil, 
		"fornax": nil,
	}
	
	type Body struct {
		CleanedBody string`json:"cleaned_body"`
	}

	cleaned_body := Body {
		CleanedBody: "",
	}

	slice := strings.SplitSeq(body, " ")
	for word := range slice{
		if _, exists := forbidden[strings.ToLower(word)]; exists {
			cleaned_body.CleanedBody += "**** "
		} else {
			cleaned_body.CleanedBody += word + " "
		}

	}
	clean := strings.Trim(cleaned_body.CleanedBody, " ")
	cleaned_body.CleanedBody = clean
	respondWithJson(w, 200, cleaned_body)
}