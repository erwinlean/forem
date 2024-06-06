package models

import "time"

type FileData struct {
    ID          string        `bson:"_id,omitempty"`
    CompanyName string        `bson:"companyName,omitempty"`
    FileName    string        `bson:"fileName,omitempty"`
    UploadedAt  time.Time     `bson:"uploadedAt,omitempty"`
    Data        interface{}   `bson:"data,omitempty"`
}