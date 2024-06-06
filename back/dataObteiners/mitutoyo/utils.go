package mitutoyo

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func writeCSV(filename string, product ProductDetail) {
	if _, exists := allScrapedProducts[product.URL]; exists {
		log.Printf("El producto ya existe en el archivo CSV: %s", product.URL)
		return
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error al abrir el archivo CSV de salida", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error al obtener información del archivo: %s", err)
		return
	}

	if fileInfo.Size() == 0 {
		header := []string{"url", "article_number", "name", "description", "short_description", "image", "technical_image", "image_links", "variants", "leaflet_links", "product_pdf_links", "instruction_pdf_links", "accesories", "youtube_links", "software_links", "attributes"}
		err := writer.Write(header)
		if err != nil {
			log.Printf("Error al escribir el encabezado del CSV: %s", err)
			return
		}
	}

	attributes := []string{}
	for key, value := range product.Attributes {
		attributes = append(attributes, fmt.Sprintf("%s: %s", key, value))
	}
	serializedAttributes := strings.Join(attributes, "; ")

	imageLinks := strings.Join(product.ImageLinks, ";")
	variants := strings.Join(product.Variants, ";")
	leafLetLinks := strings.Join(product.LeafLetLinks, ";")
	instructionPDFLinks := strings.Join(product.InstructionPDFLinks, ";")
	youtubeLinks := strings.Join(product.YoutubeLinks, ";")
	accesories := strings.Join(product.Accesories, ";")
	softwareLinks := strings.Join(product.SoftwareLinks, ";")

	record := []string{
		product.URL,
		product.ArticleNumber,
		product.Name,
		product.Description,
		product.ShortDescription,
		product.Image,
		product.TechnicalImage,
		imageLinks,
		variants,
		leafLetLinks,
		instructionPDFLinks,
		accesories,
		youtubeLinks,
		softwareLinks, 
		serializedAttributes,
	}
	err = writer.Write(record)
	if err != nil {
		log.Printf("Error al escribir los detalles del producto en el CSV: %s", err)
		return
	}

	writer.Flush()
}

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

func printProductDetails(product ProductDetail) {
	/*
	Variants: 
	LeafLet Links: 
	Instruction PDF Links: 
	Accesories: 
	*/

	//fmt.Println("Product Details:")
	//fmt.Printf("URL: %s\n", product.URL)
	//fmt.Printf("Article Number: %s\n", product.ArticleNumber)
	//fmt.Printf("Name: %s\n", product.Name)
	//fmt.Printf("Description: %s\n", product.Description)
	//fmt.Printf("Short Description: %s\n", product.ShortDescription)
	//fmt.Printf("Image: %s\n", product.Image)
	//fmt.Printf("Technical Image: %s\n", product.TechnicalImage)
	//fmt.Printf("Variants: %s\n", strings.Join(product.Variants, ", "))
	//fmt.Printf("LeafLet Links: %s\n", strings.Join(product.LeafLetLinks, ", "))
	//fmt.Printf("Instruction PDF Links: %s\n", strings.Join(product.InstructionPDFLinks, ", "))
	//fmt.Printf("Accesories: %s\n", strings.Join(product.Accesories, ", "))
	//fmt.Printf("Image Links: %s\n", strings.Join(product.ImageLinks, ", "))
	//fmt.Printf("YouTube Links: %s\n", strings.Join(product.YoutubeLinks, ", "))
	//fmt.Printf("Software Links: %s\n", strings.Join(product.SoftwareLinks, ", "))
	//fmt.Println("Attributes:")
	//for key, value := range product.Attributes {
	//	fmt.Printf("  %s: %s\n", key, value)
	//}
	//fmt.Println("------------------------------------------------")
}
