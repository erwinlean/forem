package controllers

import (
    "back/dataObteiners/mitutoyo"
    "back/utils"
    "encoding/json"
    "log"
    "net/http"
)

func Mitutoyo(w http.ResponseWriter, r *http.Request) {
    data := mitutoyo.MainMitutoyo()
    if data == nil {
        log.Println("No data obtained from Mitutoyo scraping")
        http.Error(w, "Failed to scrape products", http.StatusInternalServerError)
        return
    }

    dataJSON, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        log.Printf("Error converting data to JSON: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    recipientEmail := ""
    subject := "Scraping Mitutoyo Complete"
    body := "<html><body><div>" + string(dataJSON) + "</div></body></html>"

    err = utils.SendEmail(recipientEmail, subject, body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Printf("Error encoding data to JSON: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Println("Response sent successfully.")
}

func Cosensaws(w http.ResponseWriter, r *http.Request) {
    log.Println("in Cosensaws controller")
    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response)

    err := utils.SendEmail("recipient@example.com", "Scraping Cosensaws complete", "Scraping Cosensaws has been completed successfully!")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Cosensaws completed successfully!"))
}

func Fluke(w http.ResponseWriter, r *http.Request) {
    log.Println("in Fluke controller")
    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response)

    err := utils.SendEmail("recipient@example.com", "Scraping Fluke complete", "Scraping Fluke has been completed successfully!")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Fluke completed successfully!"))
}

func Kinkelder(w http.ResponseWriter, r *http.Request) {
    log.Println("in Kinkelder controller")
    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response)

    err := utils.SendEmail("recipient@example.com", "Scraping Kinkelder complete", "Scraping Kinkelder has been completed successfully!")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Kinkelder completed successfully!"))
}
