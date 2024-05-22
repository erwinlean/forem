package controllers

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/gorilla/sessions"
    "go.mongodb.org/mongo-driver/bson"
    "back/models"
    "back/utils"
    "golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret-key-todo"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)

    var result models.User
    err := utils.UserCollection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&result)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    session, _ := store.Get(r, "session")
    session.Values["user"] = result.Username
    session.Save(r, w)

    json.NewEncoder(w).Encode(result)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    _, err = utils.UserCollection.InsertOne(context.TODO(), user)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}