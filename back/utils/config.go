package utils

import (
    "context"
    "log"
    "os"

    "github.com/joho/godotenv"
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
    dbPassword := os.Getenv("DBpasswordAuth")
    dbHost := os.Getenv("DBhost")

    connectionString := "mongodb+srv://" + dbUser + ":" + dbPassword + "@" + dbHost + "/?retryWrites=true&w=majority"
    clientOptions := options.Client().ApplyURI(connectionString)

    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Conectado a la base de datos MongoDB")

    UserCollection = client.Database("nombre_de_tu_base_de_datos").Collection("usuarios")
}
