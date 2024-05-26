package main

import (
	"log"
	"time"
)

func main() {
	startTime := time.Now()

	mainURL := "https://shop.mitutoyo.mx/products/es_MX/index.xhtml"
	// test urls
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	// url error: https://shop.mitutoyo.mx/products/es_MX/05.05.08/Radius%20chart%20%C3%B8%20500%20mm/$catalogue/mitutoyoData/PR/63AAA548/index.xhtml
	// https://shop.mitutoyo.mx/products/es_MX/Hardness_reference_materials/40HRBW%20Rockwell%20Hardness%20Reference%20Material/$catalogue/mitutoyoData/PR/63ETB021/index.xhtml
	// https://shop.mitutoyo.mx/products/es_MX/KOMEG_Styli_Straight/Stylus%20M2%20ruby%20ball%20%C3%986%2C0mm%0A%2C/$catalogue/mitutoyoData/PR/K651543/index.xhtml
	// 
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