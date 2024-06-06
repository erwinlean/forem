package controllers

import (
    "back/dataObteiners/mitutoyo"
    "back/middleware"
    "back/utils"
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func SaveOrUpdateFileData(companyName string, data []interface{}) error {
    filter := bson.M{"companyName": companyName}
    update := bson.M{
        "$set": bson.M{
            "companyName": companyName,
            "fileName":    companyName + "_data.json",
            "uploadedAt":  time.Now(),
            "data":        data,
        },
    }
    opts := options.Update().SetUpsert(true)

    _, err := utils.FileDataCollection.UpdateOne(context.Background(), filter, update, opts)
    return err
}

func Mitutoyo(w http.ResponseWriter, r *http.Request) {
    userEmail := r.Header.Get("X-User-Email")
    token := r.Header.Get("Authorization")

    log.Println("email and token obteined in the header of the request scrapper mitutoyo")
    log.Println(userEmail)
    log.Println(token)

    _, err := middleware.ValidateJWTToken(token, jwtKey)
    if err != nil {
        log.Println("Token validation failed:", err)
        http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
        return
    }

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

    recipientEmail := userEmail
    subject := "Scraping Mitutoyo Complete"
    body := "hello from scrapper <html><body><div>" + string(dataJSON) + "</div></body></html>"

    err = utils.SendEmail(recipientEmail, subject, body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var interfaceData []interface{}
    for _, item := range data {
        interfaceData = append(interfaceData, item)
    }

    // db
    err = SaveOrUpdateFileData("Mitutoyo", interfaceData)
    if err != nil {
        log.Printf("Error saving data to database: %v", err)
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
