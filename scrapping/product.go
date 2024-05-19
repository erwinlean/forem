package main

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
