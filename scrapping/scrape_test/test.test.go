package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type ProductDetail struct {
	URL                 string
	ArticleNumber       string
	Name                string
	Description         string
	ShortDescription    string
	Image               string
	TechnicalImage      string
	Variants            []string
	LeafLetLinks        []string
	InstructionPDFLinks []string
	Accesories          []string
	ImageLinks          []string
	YoutubeLinks        []string
	SoftwareLinks       []string
	Attributes          map[string]string
}

func scrapeProduct(url string) {
	c := colly.NewCollector(
		colly.AllowedDomains("shop.mitutoyo.mx"),
		colly.Async(true),
	)

	c.OnHTML("#product", func(e *colly.HTMLElement) {
		productURL := e.Request.URL.String()
		productArticleNumber := e.ChildText("#product div.product-content-wrapper div.product-content div.articlenumber h2 span.value")
		productName := e.ChildText("#product div.product-content-wrapper div.product-content div.name h2")

		// Descripción del producto
		var productDescriptionBuilder strings.Builder
		e.ForEach("span.description", func(_ int, elem *colly.HTMLElement) {
			elem.ForEach("*", func(_ int, childElem *colly.HTMLElement) {
				productDescriptionBuilder.WriteString(childElem.Text + " ")
			})
		})
		productDescription := strings.TrimSpace(productDescriptionBuilder.String())

		productShortDescription := e.ChildText("#product div.product-content-wrapper div.product-content div.name span")
		productImage := e.ChildAttr("#product div.product-content-wrapper div.product-image-wrapper div img", "src")
		productTechnical := e.ChildAttr("#productDetailPage\\:accform\\:drawings\\:productImageColorboxLink", "href")

		// Variantes del producto
		// NO FUNCIONAL
		var variants []string
		e.ForEach("#product_variants .articlenumber a", func(_ int, elem *colly.HTMLElement) {
			log.Printf("Encontrado enlace variante: %s", elem.Attr("href"))
			variant := elem.Text
			log.Printf("Texto de variante: %s", variant)
			variants = append(variants, variant)
		})

		// Accesorios
		// NO FUNCIONAL
		var accesories []string
		e.ForEach("#productDetailPage\\:accform\\:product_relation_navi0 a", func(_ int, elem *colly.HTMLElement) {
			accesory := elem.Text
			accesories = append(accesories, accesory)
		})

		var imageLinks []string
		e.ForEach("#productDetailPage\\:thumbnails div.thumbnail a", func(_ int, elem *colly.HTMLElement) {
			imageURL := elem.Attr("href")
			absoluteImageURL := e.Request.AbsoluteURL(imageURL)
			imageLinks = append(imageLinks, absoluteImageURL)
		})
		imageLinks = imageLinks[:copy(imageLinks, imageLinks[1:])]

		// Folletos
		// NO FUNCIONAL
		var leafLetLinks []string
		e.ForEach("#product_leafLet div div div a", func(_ int, elem *colly.HTMLElement) {
			link := elem.Attr("href")
			absoluteLink := e.Request.AbsoluteURL(link)
			leafLetLinks = append(leafLetLinks, absoluteLink)
		})

		// Enlaces PDF de instrucciones
		// NO FUNCIONAL
		var instructionPDFLinks []string
		e.ForEach("#productDetailPage\\:accform\\:instructionsAcc i", func(_ int, elem *colly.HTMLElement) {
			if elem.Text == "Instrucciones" {
				e.ForEach("#productDetailPage\\:accform\\:j_idt\\d+\\:j_idt\\d+\\:0\\:downloadLink a[href$='.pdf']", func(_ int, pdfElem *colly.HTMLElement) {
					pdfLink := pdfElem.Attr("href")
					absolutePDFLink := e.Request.AbsoluteURL(pdfLink)
					instructionPDFLinks = append(instructionPDFLinks, absolutePDFLink)
				})
			}
		})

		// Enlaces de YouTube
		var youtubeLinks []string
		e.ForEach("#YoutubeVideos iframe", func(_ int, elem *colly.HTMLElement) {
			link := elem.Attr("src")
			absoluteLink := e.Request.AbsoluteURL(link)
			youtubeLinks = append(youtubeLinks, absoluteLink)
		})

		// Enlaces de software .zip
		var softwareLinks []string
		e.ForEach("#productDetailPage\\:accform\\:softwareContent a[href$='.zip']", func(_ int, elem *colly.HTMLElement) {
			link := elem.Attr("href")
			absoluteLink := e.Request.AbsoluteURL(link)
			softwareLinks = append(softwareLinks, absoluteLink)
		})

		// Atributos del producto
		attributes := make(map[string]string)
		e.ForEach("#product_parameters table tbody tr", func(_ int, elem *colly.HTMLElement) {
			param := elem.ChildText("td.parameter span span")

			var valueBuilder strings.Builder
			elem.ForEach("td.value", func(_ int, valueElem *colly.HTMLElement) {
				if valueElem.DOM.Find("ul").Length() > 0 {
					valueElem.ForEach("ul > li", func(_ int, liElem *colly.HTMLElement) {
						valueBuilder.WriteString(liElem.Text + " ")
					})
				} else {
					valueElem.ForEach("p", func(_ int, pElem *colly.HTMLElement) {
						valueBuilder.WriteString(pElem.Text + " ")
					})
					valueElem.ForEach("span", func(_ int, spanElem *colly.HTMLElement) {
						valueBuilder.WriteString(spanElem.Text + " ")
					})
					valueBuilder.WriteString(valueElem.Text + " ")
				}
			})
			value := strings.TrimSpace(valueBuilder.String())

			if param != "" && value != "" {
				attributes[param] = value
			}
		})

		product := ProductDetail{
			URL:                 productURL,
			ArticleNumber:       productArticleNumber,
			Name:                productName,
			Description:         productDescription,
			ShortDescription:    productShortDescription,
			Image:               e.Request.AbsoluteURL(productImage),
			TechnicalImage:      e.Request.AbsoluteURL(productTechnical),
			Variants:            variants,
			LeafLetLinks:        leafLetLinks,
			InstructionPDFLinks: instructionPDFLinks,
			Accesories:          accesories,
			ImageLinks:          imageLinks,
			YoutubeLinks:        youtubeLinks,
			SoftwareLinks:       softwareLinks,
			Attributes:          attributes,
		}

		log.Println("Variantes del Producto:")
		for _, variant := range product.Variants {
			log.Println("-", variant)
		}

		log.Println("Enlaces de Folletos del Producto:")
		for _, link := range product.LeafLetLinks {
			log.Println("-", link)
		}

		log.Println("Enlaces de Instrucciones PDF del Producto:")
		for _, link := range product.InstructionPDFLinks {
			log.Println("-", link)
		}

		log.Println("Accesorios del Producto:")
		for _, accessory := range product.Accesories {
			log.Println("-", accessory)
		}
	})

	c.Visit(url)
	c.Wait()
}

func main() {
	url := "https://shop.mitutoyo.mx/products/es_MX/1292249246959/Surftest%20SJ-210%20%5Bmm%5D/$catalogue/mitutoyoData/PR/178-560-11D/index.xhtml"
	scrapeProduct(url)
}
