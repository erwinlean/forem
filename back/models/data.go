package models

import "time"

type FileData struct {
    ID         string        `bson:"_id,omitempty"`
    Filename   string        `bson:"filename,omitempty"`
    UploadedAt time.Time     `bson:"uploadedAt,omitempty"`
    Data       []interface{} `bson:"data,omitempty"`
}