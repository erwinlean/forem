package main

import (
	"log"
	"time"
)

func main() {
	startTime := time.Now()

	//mainURL := "https://shop.mitutoyo.mx/products/es_MX/index.xhtml"
	mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
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