package controllers

import (
    "net/http"
    "encoding/json"
    
    //"back/utils"
)

func Mitutoyo(w http.ResponseWriter, r *http.Request) {

    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response) 

    /*
    err := utils.SendEmail("recipient@example.com", "Scraping Mitutoyo complete", "Scraping Mitutoyo has been completed successfully!")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Mitutoyo completed successfully!"))
    */
}
func Cosensaws(w http.ResponseWriter, r *http.Request) {

    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response) 

    /*
    err := utils.SendEmail("recipient@example.com", "Scraping Cosensaws complete", "Scraping Cosensaws has been completed successfully!")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Cosensaws completed successfully!"))
    */
}
func Fluke(w http.ResponseWriter, r *http.Request) {

    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response) 

    /*
    err := utils.SendEmail("recipient@example.com", "Scraping Fluke complete", "Scraping Fluke has been completed successfully!")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Fluke completed successfully!"))
    */
}
func Kinkelder(w http.ResponseWriter, r *http.Request) {

    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response) 

    /*
    err := utils.SendEmail("recipient@example.com", "Scraping Kinkelder complete", "Scraping Kinkelder has been completed successfully!")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Kinkelder completed successfully!"))
    */
}
