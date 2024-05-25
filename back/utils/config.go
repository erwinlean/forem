package utils

import (
    "context"
    "log"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
    clientOptions := options.Client().ApplyURI("mongo_uri")

    // Crea el cliente con las opciones configuradas
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Verifica si la conexión se estableció correctamente
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB!")
    return client
}

var DB *mongo.Client = ConnectDB()
var UserCollection *mongo.Collection = DB.Database("todo").Collection("users")