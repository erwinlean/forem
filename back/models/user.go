package models

type User struct {
    ID       string `bson:"_id,omitempty"`
    Username string `bson:"username,omitempty"`
    Password string `bson:"password,omitempty"`
}