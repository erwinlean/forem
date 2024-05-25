package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

var urlCount int = 0
var productCount int = 0
var visitedURLs map[string]bool = make(map[string]bool)
var visitedImages map[string]bool = make(map[string]bool)

type ProductInfo struct {
	Src            string
	Link           string
	ShortDescription string
	Description    string
}

func main() {
	startTime := time.Now()

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}

	baseURL := "https://www.cosensaws.com/saws"
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	if _, err = page.Goto(baseURL); err != nil {
		log.Fatalf("could not go to base URL: %v", err)
	}

	time.Sleep(2 * time.Second) // Wait for 2 seconds to let the page load

	initialLinks, err := extractLinks(page, "a[data-testid='linkElement']")
	if err != nil {
		log.Fatalf("could not extract initial links: %v", err)
	}

	file, err := os.Create("product_images.csv")
	if err != nil {
		log.Fatalf("could not create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Product URL", "Image URL", "Short Description", "Description"})

	for _, link := range initialLinks {
		if !visitedURLs[link] {
			visitedURLs[link] = true
			nestedLinks, err := getNestedLinks(context, link, "a[data-testid='linkElement']")
			if err != nil {
				log.Printf("could not get nested links from %s: %v", link, err)
				continue
			}

			nestedLinks = filterLinks(nestedLinks)
			nestedLinks = eliminar(nestedLinks)

			for _, nestedLink := range nestedLinks {
				if !visitedURLs[nestedLink] {
					visitedURLs[nestedLink] = true
					productInfo, err := getProductInfo(context, nestedLink)
					if err != nil {
						log.Printf("could not get product info from %s: %v", nestedLink, err)
						continue
					}
					writeProductInfo(writer, productInfo)

					if !visitedImages[productInfo.Src] {
						visitedImages[productInfo.Src] = true
						writer.Write([]string{nestedLink, productInfo.Src, productInfo.ShortDescription, productInfo.Description})
						fmt.Printf("Src: %s, Link: %s, Short Description: %s, Description: %s\n", productInfo.Src, nestedLink, productInfo.ShortDescription, productInfo.Description)
					}
				}
			}
		}
	}

	log.Printf("Urls processed: %d", urlCount)
	log.Printf("Products processed: %d", productCount)
	// Check time of the scraping
	duration := time.Since(startTime)
	log.Printf("Scraping completed in %v\n", duration)
}

func extractLinks(page playwright.Page, selector string) ([]string, error) {
	elements, err := page.QuerySelectorAll(selector)
	if err != nil {
		return nil, err
	}

	var links []string
	for _, element := range elements {
		href, err := element.GetAttribute("href")
		if err != nil {
			return nil, err
		}
		if href != "" {
			urlCount++
			log.Printf("URL Count: %d, URL: %s", urlCount, href)
			links = append(links, href)
		}
	}

	return links, nil
}

func getNestedLinks(context playwright.BrowserContext, url, selector string) ([]string, error) {
	urlCount++
	log.Printf("URL Count: %d, URL: %s", urlCount, url)

	page, err := context.NewPage()
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if _, err = page.Goto(url); err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Second) // Wait for 2 seconds to let the page load

	for {
		previousHeight, _ := page.Evaluate("document.body.scrollHeight")
		page.Evaluate("window.scrollTo(0, document.body.scrollHeight)")
		time.Sleep(4 * time.Second)
		newHeight, _ := page.Evaluate("document.body.scrollHeight")
		if previousHeight == newHeight {
			break
		}
	}

	return extractLinks(page, selector)
}

func filterLinks(links []string) []string {
	var filteredLinks []string
	re := regexp.MustCompile(`/[^/]*\d{3}[^/]*$`)

	for _, link := range links {
		if re.MatchString(link) {
			filteredLinks = append(filteredLinks, link)
		}
	}

	return filteredLinks
}

func eliminar(links []string) []string {
	var result []string
	for _, link := range links {
		if strings.HasPrefix(link, "https://www.cosensaws.com/") && !strings.Contains(link, "-saws-") && !strings.Contains(link, "southtec2023") {
			if !strings.Contains(link, "/") || strings.Count(link, "/") <= 4 {
				productCount++
				log.Printf("Product Count: %d, Product URL: %s", productCount, link)
				result = append(result, link)
			}
		}
	}
	return result
}

func getProductInfo(context playwright.BrowserContext, url string) (ProductInfo, error) {
	page, err := context.NewPage()
	if err != nil {
		return ProductInfo{}, err
	}
	defer page.Close()

	if _, err = page.Goto(url); err != nil {
		return ProductInfo{}, err
	}

	time.Sleep(4 * time.Second) // Wait for 2 seconds to let the page load

	src, err := page.EvalOnSelector("#img_undefined > img", "e => e.src", nil)
	if err != nil {
		return ProductInfo{}, err
	}

	srcStr, ok := src.(string)
	if !ok {
		return ProductInfo{}, fmt.Errorf("failed to convert src to string")
	}

	// Eliminamos el segmento /v1/fill y lo que sigue después del identificador de la imagen
	re := regexp.MustCompile(`(\.png|\.jpg|\.jpeg|\.webp)/v1/fill.*`)
	srcStr = re.ReplaceAllString(srcStr, `$1`)

	// Extraer la descripción corta
	shortDesc, err := page.EvalOnSelector(`[id^="comp-"] > h4 > span > span > span > span > span`, "e => e.textContent", nil)
	if err != nil {
	    shortDesc = ""
	}
	shortDescStr, ok := shortDesc.(string)
	if !ok {
	    shortDescStr = ""
	}

	// Extraer la descripción
	desc, err := page.EvalOnSelector(`[id^="comp-"] > p > span > span > span`, "e => e.textContent", nil)
	if err != nil {
	    desc = ""
	}
	descStr, ok := desc.(string)
	if !ok {
	    descStr = ""
	}

	return ProductInfo{
		Src:            srcStr,
		Link:           url,
		ShortDescription: shortDescStr,
		Description:    descStr,
	}, nil
}

func writeProductInfo(writer *csv.Writer, productInfo ProductInfo) {
    // Write header if the file is empty
    file, _ := os.Stat("product_images.csv")
    if file.Size() == 0 {
        writer.Write([]string{"Product URL", "Image URL", "Short Description", "Description"})
    }
    writer.Write([]string{productInfo.Link, productInfo.Src, productInfo.ShortDescription, productInfo.Description})
}