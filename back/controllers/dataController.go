package controllers

import (
    "encoding/json"
    "net/http"
)

func MitutoyoData(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Hello from MitutoyoData!"}
    json.NewEncoder(w).Encode(response)
}

func CosensawsData(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Hello from CosensawsData!"}
    json.NewEncoder(w).Encode(response)
}

func FlukeData(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Hello from FlukeData!"}
    json.NewEncoder(w).Encode(response)
}

func KinkelderData(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Hello from KinkelderData!"}
    json.NewEncoder(w).Encode(response)
}
