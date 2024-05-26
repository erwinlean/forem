package main

import (
	"log"
	"time"
)

func main() {
	startTime := time.Now()

	mainURL := "https://shop.mitutoyo.mx/products/es_MX/index.xhtml"
	// Test urls 4 missing attr
	// Table data
	// table nueva de atributos: 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/KOMEG_Styli_Straight/Stylus%20M2%20ruby%20ball%20%C3%986%2C0mm%0A%2C/$catalogue/mitutoyoData/PR/K651543/index.xhtml"
	// variantes: 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	// accesorios: 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	// "Software": 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	// "Folleto, instrucciones, CAD": 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/01.02.07.01/164-163/$catalogue/mitutoyoData/PR/164-163/index.xhtml"

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

	log.Print("Url scrapped: ", urlCounts)
}