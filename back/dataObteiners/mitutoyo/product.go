package mitutoyo

// ProductDetail represents the details of a product
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