package controllers

import (
    "back/dataObteiners/mitutoyo"
    "back/middleware"
    "back/utils"
    "back/models"

    "os"
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

    var user models.User
    filter := bson.M{"email": userEmail}
    err := utils.UserCollection.FindOne(context.Background(), filter).Decode(&user)
    if err != nil {
        log.Println("User not found:", err)
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    _, err = middleware.ValidateJWTToken(token, jwtKey)
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

    // test csv
    // 1- ver funcionamiento correcto funcion EmailStyle > funciono, mejorar ui.
    // 2- ver que funciona la creacion csv > creado el csv de forma correct en back/products.csv (main path)
    // 3- ver envio atachment de csv
    csvFile := "products.csv"
	for i := 0; i < len(data); i++ {
		log.Println(data[i])
		mitutoyo.WriteCSV(csvFile, data[i])
	}   

    dataJSON, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        log.Printf("Error converting data to JSON: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    recipientEmail := userEmail
    subject := "Scraping Mitutoyo Complete user " + user.Username
    body := utils.EmailStyle(user.Username, "Mitutoyo", string(dataJSON))

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

    // remove csv
    err = os.Remove(csvFile)
	if err != nil {
		log.Printf("Error deleting the CSV file: %v", err)
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
