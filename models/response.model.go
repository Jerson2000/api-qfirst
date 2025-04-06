package models

import (
	"encoding/json"
	"net/http"
)

type AuthResponse struct {
	Message string `json:"message"`
	UserId  uint   `json:"userId"`
	Name    string `json:"name"`
	Token   string `json:"token"`
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
