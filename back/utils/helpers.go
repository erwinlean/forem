package utils

import (
    "encoding/json"
    "net/http"
)

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

func RespondWithError(w http.ResponseWriter, status int, message string) {
    RespondWithJSON(w, status, map[string]string{"error": message})
}