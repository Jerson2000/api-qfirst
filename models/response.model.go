package models

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// A response when logging,signing up and reissuance of a token
// @Schema
type AuthResponse struct {
	// Respresents the message of the response
	// @example Signup successful
	Message string `json:"message"`
	// Represent the id of the user
	UserId uuid.UUID `json:"userId"`
	// Represent the name of the user
	Name string `json:"name"`
	// Represent the access token for the user
	Token string `json:"token"`
}

func ResponseWithError(res http.ResponseWriter, code int, message string) {
	ResponseWithJSON(res, code, map[string]string{"error": message})
}

func ResponseWithJSON(res http.ResponseWriter, code int, payload any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	if err := json.NewEncoder(res).Encode(payload); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
