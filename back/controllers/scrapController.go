package controllers

import (
	"back/dataObteiners/mitutoyo"
	"back/middleware"
	"back/models"
	"back/utils"
	//"reflect"
	//"strings"

	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveOrUpdateFileData(companyName string, data []mitutoyo.ProductDetail) error {
    filter := bson.M{"company": companyName}
    update := bson.M{
        "$set": bson.M{
            "company":    companyName,
            "fileName":   companyName + "_data.json",
            "uploadedAt": time.Now(),
            "products":   convertToProductStructs(data),
        },
    }
    opts := options.Update().SetUpsert(true)

    _, err := utils.FileDataCollection.UpdateOne(context.Background(), filter, update, opts)
    if err != nil {
        log.Printf("Error updating database: %v", err)
        return err
    }

    return nil
}

// Función para convertir []mitutoyo.ProductDetail a []models.Product
func convertToProductStructs(data []mitutoyo.ProductDetail) []models.Product {
	var products []models.Product

	for _, detail := range data {
		product := models.Product{
			URL:                 detail.URL,
			ArticleNumber:       detail.ArticleNumber,
			Name:                detail.Name,
			Description:         detail.Description,
			ShortDescription:    detail.ShortDescription,
			Image:               detail.Image,
			TechnicalImage:      detail.TechnicalImage,
			Variants:            detail.Variants,
			LeafLetLinks:        detail.LeafLetLinks,
			InstructionPDFLinks: detail.InstructionPDFLinks,
			Accesories:          detail.Accesories,
			ImageLinks:          detail.ImageLinks,
			YoutubeLinks:        detail.YoutubeLinks,
			SoftwareLinks:       detail.SoftwareLinks,
			Attributes:          detail.Attributes,
		}
		products = append(products, product)
	}

	return products
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

    dbData, err := fetchFileData("Mitutoyo")
    if err != nil {
        log.Printf("Error fetching data from database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // getting != betwen new and old (based on the db) then re-writes the db with all of them
    // sending the email by 2 csv, with newest products and all of them
    var allData, newData []mitutoyo.ProductDetail
    existingProducts := make(map[string]bool)
    for _, productsDb := range dbData {
        for _, productDb := range productsDb.Products {
            existingProducts[productDb.ArticleNumber] = true
            allData = append(allData, mitutoyo.ProductDetail{
                URL:                 productDb.URL,
                ArticleNumber:       productDb.ArticleNumber,
                Name:                productDb.Name,
                Description:         productDb.Description,
                ShortDescription:    productDb.ShortDescription,
                Image:               productDb.Image,
                TechnicalImage:      productDb.TechnicalImage,
                Variants:            productDb.Variants,
                LeafLetLinks:        productDb.LeafLetLinks,
                InstructionPDFLinks: productDb.InstructionPDFLinks,
                Accesories:          productDb.Accesories,
                ImageLinks:          productDb.ImageLinks,
                YoutubeLinks:        productDb.YoutubeLinks,
                SoftwareLinks:       productDb.SoftwareLinks,
                Attributes:          productDb.Attributes,
            })
        }
    }

    for _, product := range data {
        if !existingProducts[product.ArticleNumber] {
            newData = append(newData, product)
        }
        allData = append(allData, product)
    }

    csvFile := "mitutoyo_complete_data.csv"
    newCsvFile := "mitutoyo_new_data.csv"

    for _, product := range allData {
        mitutoyo.WriteCSV(csvFile, 0, product)
    }
    for _, product := range newData {
        mitutoyo.WriteCSV(newCsvFile, 1, product)
    }

    // Enviar email
    recipientEmail := userEmail
    subject := "Scraping Mitutoyo Completada " + user.Username
    body := utils.EmailStyle(user.Username, "Mitutoyo")
    err = utils.SendEmail(recipientEmail, subject, body, csvFile, newCsvFile)
    if err != nil {
        log.Printf("Error sending email: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Guardar en la db
    err = SaveOrUpdateFileData("Mitutoyo", allData)
    if err != nil {
        log.Printf("Error saving data to database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Limpiar archivos CSV
    err = os.Truncate(csvFile, 0)
    if err != nil {
        log.Printf("Error clearing the CSV file: %v", err)
    }
    err = os.Truncate(newCsvFile, 0)
    if err != nil {
        log.Printf("Error clearing the CSV file: %v", err)
    }
    mitutoyo.ResetScrappedProducts()

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(allData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("Response sent successfully.")
}

func Cosensaws(w http.ResponseWriter, r *http.Request) {
    log.Println("in Cosensaws controller")
    response := map[string]string{"message": "Hello test!"}
    json.NewEncoder(w).Encode(response)

    err := utils.SendEmail("recipient@example.com", "Scraping Cosensaws complete", "Scraping Cosensaws has been completed successfully!","", "")
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

    err := utils.SendEmail("recipient@example.com", "Scraping Fluke complete", "Scraping Fluke has been completed successfully!","", "")
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

    err := utils.SendEmail("recipient@example.com", "Scraping Kinkelder complete", "Scraping Kinkelder has been completed successfully!","", "")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Scraping Kinkelder completed successfully!"))
}