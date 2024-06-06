package utils

import (
    "encoding/json"
    "net/http"
)

// RespondWithJSON responde con un objeto JSON y un estado HTTP.
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(response)
}

// RespondWithError responde con un mensaje de error y un estado HTTP.
func RespondWithError(w http.ResponseWriter, status int, message string) {
    RespondWithJSON(w, status, map[string]string{"error": message})
}