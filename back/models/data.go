package models

import "time"

type ProductData struct {
	ID          string    `bson:"_id,omitempty"`
	Company     string    `bson:"company,omitempty"`
	FileName    string    `bson:"fileName,omitempty"`
	UploadedAt  time.Time `bson:"uploadedAt,omitempty"`
	Products    []Product `bson:"products,omitempty"`
}

type Product struct {
	URL                 string            `bson:"url,omitempty"`
	ArticleNumber       string            `bson:"articleNumber,omitempty"`
	Name                string            `bson:"name,omitempty"`
	Description         string            `bson:"description,omitempty"`
	ShortDescription    string            `bson:"shortDescription,omitempty"`
	Image               string            `bson:"image,omitempty"`
	TechnicalImage      string            `bson:"technicalImage,omitempty"`
	Variants            []string          `bson:"variants,omitempty"`
	LeafLetLinks        []string          `bson:"leafLetLinks,omitempty"`
	InstructionPDFLinks []string          `bson:"instructionPDFLinks,omitempty"`
	Accesories          []string          `bson:"accesories,omitempty"`
	ImageLinks          []string          `bson:"imageLinks,omitempty"`
	YoutubeLinks        []string          `bson:"youtubeLinks,omitempty"`
	SoftwareLinks       []string          `bson:"softwareLinks,omitempty"`
	Attributes          map[string]string `bson:"attributes,omitempty"`
}