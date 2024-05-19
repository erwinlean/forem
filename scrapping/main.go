package main

import (
	"log"
	"time"
)

func main() {
	startTime := time.Now()

	mainURL := "https://shop.mitutoyo.mx/products/es_MX/index.xhtml"
	ok := checkURL(mainURL)
	if !ok {
		log.Println("Skipping scraping as mainURL is not accessible.")
		return
	}

	// Init
	scrapeRecursive(mainURL)

	// Check time of the scraping
	duration := time.Since(startTime)
	log.Printf("Scraping completed in %v\n", duration)
}
