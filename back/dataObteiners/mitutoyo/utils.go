package mitutoyo

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

var allScrappedProducts = make(map[string]bool)

func ResetScrappedProducts() {
    allScrappedProducts = make(map[string]bool)
}

func WriteCSV(filename string, fileNumber int ,product ProductDetail) {

	if(fileNumber == 0){
		if _, exists := allScrappedProducts[product.URL]; exists {
			log.Printf("El producto ya existe en el archivo CSV: %s", product.URL)
			return
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error al abrir el archivo CSV de salida", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

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

	allScrappedProducts[product.URL] = true
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
