package utils

import (
    "context"
    "log"
    "os"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    
    dbUser := os.Getenv("DBuser")
    dbPassword := os.Getenv("DBpassword")
    dbHost := os.Getenv("DBhost")

    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI("mongodb+srv://" + dbUser + ":" + dbPassword + "@" + dbHost + "/?retryWrites=true&w=majority&appName=foremDB").SetServerAPIOptions(serverAPI)
    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            panic(err)
        }
    }()

    if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
        panic(err)
    }

    log.Println("Conectado a la base de datos MongoDB")
    UserCollection = client.Database("foremDB").Collection("users")
}