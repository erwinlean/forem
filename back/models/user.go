package models

import "time"

type User struct {
    ID         string      `bson:"_id,omitempty"`
    Username   string      `bson:"username,omitempty"`
    Password   string      `bson:"password,omitempty"`
    Email      string      `bson:"email,omitempty"`
    CreatedAt  time.Time   `bson:"createdAt,omitempty"`
    LoginDates []time.Time `bson:"loginDates,omitempty"`
}