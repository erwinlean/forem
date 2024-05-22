package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

var (
	allScrapedURLs     = make(map[string]bool)
	allScrapedProducts = make(map[string]bool)
)

func scrapeRecursive(url string) {
	if _, visited := allScrapedURLs[url]; visited {
		return
	}

	subURLs, products := scrapeCategories(url)
	allScrapedURLs[url] = true

	for _, subURL := range subURLs {
		scrapeRecursive(subURL)
	}

	for _, product := range products {
		allScrapedProducts[product.URL] = true
	}
}

func scrapeCategories(url string) ([]string, []ProductDetail) {
	c := colly.NewCollector(
		colly.AllowedDomains("shop.mitutoyo.mx"),
		colly.Async(false), // Cambiar a false para pruebas
	)

	var subURLs []string
	var products []ProductDetail

	c.OnHTML("form#categoryWrapper\\:form a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("id"), "categoryWrapper") {
			subURL := e.Request.AbsoluteURL(e.Attr("href"))
			subURLs = append(subURLs, subURL)
		}
	})

	c.OnHTML("table#categoryWrapper\\:categoryObjectsTable\\:CT_\\:table a[href]", func(e *colly.HTMLElement) {
		subURL := e.Request.AbsoluteURL(e.Attr("href"))
		subURLs = append(subURLs, subURL)
	})

	c.OnHTML("table#productFamily\\:productFamilyObjectsTable\\:PG\\:table a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("id"), "productFamily") {
			subURL := e.Request.AbsoluteURL(e.Attr("href"))
			subURLs = append(subURLs, subURL)
		}
	})

	c.OnHTML("td.dataScroller_pageSelector ul li.pageSelector_item a[href]", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Attr("class"), "pageSelector_disabled") {
			nextPageURL := e.Request.AbsoluteURL(e.Attr("href"))
			c.Visit(nextPageURL)
		}
	})

	// Product data
	c.OnHTML("#product", func(e *colly.HTMLElement) {
		product := extractProductDetails(e)
		writeCSV("products.csv", product)
		products = append(products, product)

		printProductDetails(product)
	})

	c.Visit(url)
	c.Wait()
	
	return subURLs, products
}

func extractProductDetails(e *colly.HTMLElement) ProductDetail {
	productURL := e.Request.URL.String()
	productArticleNumber := e.ChildText("#product div.product-content-wrapper div.product-content div.articlenumber h2 span.value")
	productName := e.ChildText("#product div.product-content-wrapper div.product-content div.name h2")
	
	// Descripción del producto
	var productDescriptionBuilder strings.Builder
	e.ForEach(".description", func(_ int, elem *colly.HTMLElement) {
		elem.DOM.Contents().Each(func(_ int, node *goquery.Selection) {
			productDescriptionBuilder.WriteString(strings.TrimSpace(node.Text()) + " ")
		})
	})
	productDescription := strings.TrimSpace(productDescriptionBuilder.String())
	
	productShortDescription := e.ChildText("#product div.product-content-wrapper div.product-content div.name span")
	productImage := e.ChildAttr("#product div.product-content-wrapper div.product-image-wrapper div img", "src")
	if strings.Contains(productImage, "image_not_available_web") {
		productImage = ""
	}else{
		productImage = e.Request.AbsoluteURL(productImage)
	}
	productTechnical := e.ChildAttr("#productDetailPage\\:accform\\:drawings\\:productImageColorboxLink", "href")
	if productTechnical != ""{
		productTechnical = e.Request.AbsoluteURL(productTechnical)
	}
	
	var imageLinks []string
	e.ForEach("#productDetailPage\\:thumbnails div.thumbnail a", func(_ int, elem *colly.HTMLElement) {
		imageURL := elem.Attr("href")
		if !strings.Contains(imageURL, "image_not_available_web") {
			if imageURL != "" {
				absoluteImageURL := e.Request.AbsoluteURL(imageURL)
				imageLinks = append(imageLinks, absoluteImageURL)
			}
		}
	})
	if len(imageLinks) > 0 && imageLinks[0] == productImage {
		imageLinks = imageLinks[1:]
	}
	
	// Variantes del producto > product_variants
	// NO FUNCIONAL #tileWrapper > div > div > div:nth-child(1) > div > div > div.general > div.articlenumber > a	
	// testing width chromedp for click on x functions check onclick="..."
	var variants []string
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var variantNodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(productURL),
		chromedp.WaitVisible("#productDetailPage\\:accform\\:variantsAcc > i"),
		chromedp.Click("#productDetailPage\\:accform\\:variantsAcc > i"),
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes("a.listLink", &variantNodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Println("Error al extraer las variantes:", err)
	}

	for _, node := range variantNodes {
		variants = append(variants, node.NodeValue)
	}
	log.Print(variants)

	// Accesorios > product_relations
	// NO FUNCIONAL
	var accesories []string
	e.ForEach("#productDetailPage\\:accform\\:product_relation_navi0 a", func(_ int, elem *colly.HTMLElement) {
		accesory := elem.Text
		accesories = append(accesories, accesory)
	})

	// Folletos
	// NO FUNCIONAL
	var leafLetLinks []string
	e.ForEach("#product_leafLet div div div a", func(_ int, elem *colly.HTMLElement) {
		link := elem.Attr("href")
		if link != "" {
			absoluteLink := e.Request.AbsoluteURL(link)
			leafLetLinks = append(leafLetLinks, absoluteLink)
		}
	})
	
	// Enlaces PDF de instrucciones
	// NO FUNCIONAL
	var instructionPDFLinks []string
	e.ForEach("#productDetailPage\\:accform\\:instructionsAcc i", func(_ int, elem *colly.HTMLElement) {
		if elem.Text == "Instrucciones" {
			e.ForEach("#productDetailPage\\:accform\\:j_idt\\d+\\:j_idt\\d+\\:0\\:downloadLink a[href$='.pdf']", func(_ int, pdfElem *colly.HTMLElement) {
				pdfLink := pdfElem.Attr("href")
				if pdfLink != "" {
					absolutePDFLink := e.Request.AbsoluteURL(pdfLink)
					instructionPDFLinks = append(instructionPDFLinks, absolutePDFLink)
				}
			})
		}
	})
	
	// Enlaces de YouTube
	var youtubeLinks []string
	e.ForEach("#YoutubeVideos iframe", func(_ int, elem *colly.HTMLElement) {
		link := elem.Attr("src")
		if link != "" {
			absoluteLink := e.Request.AbsoluteURL(link)
			youtubeLinks = append(youtubeLinks, absoluteLink)
		}
	})
	
	// Enlaces de software .zip
	var softwareLinks []string
	e.ForEach("#productDetailPage\\:accform\\:softwareContent a[href$='.zip']", func(_ int, elem *colly.HTMLElement) {
		link := elem.Attr("href")
		if link != "" {
			absoluteLink := e.Request.AbsoluteURL(link)
			softwareLinks = append(softwareLinks, absoluteLink)
		}
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

	return ProductDetail {
		URL:                 productURL,
		ArticleNumber:       productArticleNumber,
		Name:                productName,
		Description:         productDescription,
		ShortDescription:    productShortDescription,
		Image:               productImage,
		TechnicalImage:      productTechnical,
		Variants:            variants,
		LeafLetLinks:        leafLetLinks,
		InstructionPDFLinks: instructionPDFLinks,
		Accesories:          accesories,
		ImageLinks:          imageLinks,
		YoutubeLinks:        youtubeLinks,
		SoftwareLinks:       softwareLinks,
		Attributes:          attributes,
	}
}