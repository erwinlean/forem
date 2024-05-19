package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Create CSV (hardcoded) of the data
func writeCSV(filename string, product ProductDetail) {
	if _, exists := allScrapedProducts[product.URL]; exists {
		log.Printf("Product already exists in CSV: %s", product.URL)
		return
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Verify if exists
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Failed to get file info: %s", err)
		return
	}

	// Header
	if fileInfo.Size() == 0 {
		header := []string{"url", "article_number", "name", "description", "short_description", "image", "technical_image", "information_pdf", "spec_sheet_mitutoyo", "image_links", "attributes"}
		err := writer.Write(header)
		if err != nil {
			log.Printf("Failed to write CSV header: %s", err)
			return
		}
	}

	// Serialize attributes map to a single string
	attributes := []string{}
	for key, value := range product.Attributes {
		attributes = append(attributes, fmt.Sprintf("%s: %s", key, value))
	}
	serializedAttributes := strings.Join(attributes, "; ")

	// Write CSV
	imageLinks := strings.Join(product.ImageLinks, ";")
	record := []string{
		product.URL,
		product.ArticleNumber,
		product.Name,
		product.Description,
		product.ShotDescription,
		product.Image,
		product.TechnicalImage,
		product.InformationPDF,
		product.SpecSheetMitutoyo,
		imageLinks,
		serializedAttributes,
	}
	err = writer.Write(record)
	if err != nil {
		log.Printf("Failed to write product details to CSV: %s", err)
		return
	}

	writer.Flush()
}

// Check if the URL is working
func checkURL(url string) bool {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Successfully fetched %s with status code %d", url, resp.StatusCode)
		return true
	}

	log.Printf("Failed to fetch %s with status code %d", url, resp.StatusCode)
	return false
}
