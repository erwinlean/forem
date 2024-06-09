package mitutoyo

import (
	"log"
	"time"
)

func MainMitutoyo() []ProductDetail { // Changed the return type to []ProductDetail
	startTime := time.Now()

	//mainURL := "https://shop.mitutoyo.mx/products/es_MX/index.xhtml"
	// Test urls 4 missing attr
	// Table data
	// table nueva de atributos: 
	//mainURL := "https://shop.mitutoyo.mx/products/es_MX/KOMEG_Styli_Straight/Stylus%20M2%20ruby%20ball%20%C3%986%2C0mm%0A%2C/$catalogue/mitutoyoData/PR/K651543/index.xhtml"
	// variantes: 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	// accesorios: 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	// "Software": 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	// "Folleto, instrucciones, CAD": 
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/01.02.07.01/164-163/$catalogue/mitutoyoData/PR/164-163/index.xhtml"
	// prueba url de 101 productos:
	// mainURL := "https://shop.mitutoyo.mx/products/es_MX/1298540172432/Sensor%20Systems/index.xhtml"
	// second test urls
	mainURL := "https://shop.mitutoyo.mx/products/es_MX/03/Micr%C3%B3metros%20Laser/index.xhtml;jsessionid=55B66B4B2EDB0B0B0D20C5106D5714CF"

	ok := checkURL(mainURL)
	if !ok {
		log.Println("Skipping scraping as mainURL is not accessible.")
		return nil // Return nil if the URL is not accessible
	}

	// Init
	scrapeRecursive(mainURL)

	// Check time of the scraping
	duration := time.Since(startTime)
	log.Printf("Scraping completed in %v\n", duration)

	log.Print("URL scrapped: ", urlCounts)

	// Convert map to slice
	var products []ProductDetail
	for _, product := range allScrapedProducts {
		products = append(products, product)
	}

	return products
}