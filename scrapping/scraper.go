package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// Check for repeated
var allScrapedURLs = make(map[string]bool)
var allScrapedProducts = make(map[string]bool)

func scrapeRecursive(url string) {
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
