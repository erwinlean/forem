package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// ProductDetail represents the details of a product
type ProductDetail struct {
	URL                string
	ArticleNumber      string
	Name               string
	Description        string
	ShotDescription    string
	Image              string
	TechnicalImage     string
	InformationPDF     string
	SpecSheetMitutoyo  string
	ImageLinks         []string
	Attributes         map[string]string
}

// Check for repeated
var allScrapedURLs = make(map[string]bool)
var allScrapedProducts = make(map[string]bool)

func main() {
	startTime := time.Now()

	mainURL := "https://shop.mitutoyo.mx/products/es_MX/index.xhtml"
	ok := checkURL(mainURL)
	if !ok {
		log.Println("Skipping scraping as mainURL is not accessible.")
		return
	}

	// Recursive urls 
	var scrapeRecursive func(url string)
	scrapeRecursive = func(url string) {
		// Check if the URL has already been visited
		if _, visited := allScrapedURLs[url]; visited {
			return
		}

		// Scrape
		subURLs, products := scrapeCategories(url)

		// Add URL to list of collected URLs
		allScrapedURLs[url] = true

		// Sub URLs
		for _, subURL := range subURLs {
			scrapeRecursive(subURL)
		}

		// Add products to the list
		for _, product := range products {
			allScrapedProducts[product.URL] = true
		}
	}

	// Initialize
	scrapeRecursive(mainURL)

	// Check time of the scraping
	duration := time.Since(startTime)
	fmt.Printf("Scraping completed in %v\n", duration)
}

// Scrape for data
func scrapeCategories(url string) ([]string, []ProductDetail) {
	c := colly.NewCollector(
		colly.AllowedDomains("shop.mitutoyo.mx"),
		colly.Async(true),
	)

	var subURLs []string
	var products []ProductDetail

	// Categories URLs if they are in form
	c.OnHTML("form#categoryWrapper\\:form a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("id"), "categoryWrapper") {
			subURL := e.Request.AbsoluteURL(e.Attr("href"))
			subURLs = append(subURLs, subURL)
		}
	})

	// Categories URLs if they are in tables
	c.OnHTML("table#categoryWrapper\\:categoryObjectsTable\\:CT_\\:table a[href]", func(e *colly.HTMLElement) {
		subURL := e.Request.AbsoluteURL(e.Attr("href"))
		subURLs = append(subURLs, subURL)
	})

	// Search URLs from family to product URLs
	c.OnHTML("table#productFamily\\:productFamilyObjectsTable\\:PG\\:table a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("id"), "productFamily") {
			subURL := e.Request.AbsoluteURL(e.Attr("href"))
			subURLs = append(subURLs, subURL)
		}
	})

	// Pagination
	c.OnHTML("td.dataScroller_pageSelector ul li.pageSelector_item a[href]", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Attr("class"), "pageSelector_disabled") {
			nextPageURL := e.Request.AbsoluteURL(e.Attr("href"))
			c.Visit(nextPageURL)
		}
	})

	// Get data
	c.OnHTML("#product", func(e *colly.HTMLElement) {

		productURL := e.Request.URL.String()
		productArticleNumber := e.ChildText("#product div.product-content-wrapper div.product-content div.articlenumber h2 span.value")
		productName := e.ChildText("#product div.product-content-wrapper div.product-content div.name h2")
        // To improve, some  descripcion are not span, are strong, are paragraph, etc to CHECK 
		productDescription := e.ChildText("#productDetailPage > span.description")
		productShortDescription := e.ChildText("#product div.product-content-wrapper div.product-content div.name span")
		productImage := e.ChildAttr("#product div.product-content-wrapper div.product-image-wrapper div img", "src")

		// To test if is it working or not
		productTechnical := e.ChildAttr("#productDetailPage\\:accform\\:drawings\\:productImageColorboxLink", "src")
		productInformationPDF := e.ChildAttr("#productDetailPage\\:accform\\:downloadLink", "href")
		productSpecSheetMitutoyo := e.ChildAttr("#productDetailPage\\:", "href")

        // working get severals imagees of the product
		var imageLinks []string
		e.ForEach("#productDetailPage\\:thumbnails div.thumbnail a", func(_ int, elem *colly.HTMLElement) {
			imageURL := elem.Attr("href")
			absoluteImageURL := e.Request.AbsoluteURL(imageURL)
			imageLinks = append(imageLinks, absoluteImageURL)
		})

		// Getting the data, but not correctly on the CSV. To improve
		attributes := make(map[string]string)
		e.ForEach("#product_parameters table tbody tr", func(_ int, elem *colly.HTMLElement) {
			param := elem.ChildText("td.parameter span span")
			value := elem.ChildText("td.value div")
			if param != "" && value != "" {
				attributes[param] = value
			}
		})

		log.Print(productURL)

		if productName != "" && productImage != "" {
			absoluteImageURL := e.Request.AbsoluteURL(productImage)
			absoluteTechnicalImageUrl := e.Request.AbsoluteURL(productTechnical)
			absoluteInformationPDF := e.Request.AbsoluteURL(productInformationPDF)
			absoluteSpecSheetMitutoyo := e.Request.AbsoluteURL(productSpecSheetMitutoyo)

			product := ProductDetail{
				URL:                productURL,
				ArticleNumber:      productArticleNumber,
				Name:               productName,
				Description:        productDescription,
				ShotDescription:    productShortDescription,
				Image:              absoluteImageURL,
				TechnicalImage:     absoluteTechnicalImageUrl,
				InformationPDF:     absoluteInformationPDF,
				SpecSheetMitutoyo:  absoluteSpecSheetMitutoyo,
				ImageLinks:         imageLinks,
				Attributes:         attributes,
			}
			writeCSV("products.csv", product)
			products = append(products, product)
		}
	})

	c.Visit(url)
	c.Wait()

	return subURLs, products
}

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
