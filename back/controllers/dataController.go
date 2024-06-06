package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"back/models"
	"back/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func fetchFileData(companyName string) ([]models.FileData, error) {
    var results []models.FileData
    filter := bson.M{"companyName": companyName}
    cursor, err := utils.FileDataCollection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var fileData models.FileData
        if err := cursor.Decode(&fileData); err != nil {
            return nil, err
        }
        results = append(results, fileData)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return results, nil
}

func MitutoyoData(w http.ResponseWriter, r *http.Request) {
    data, err := fetchFileData("Mitutoyo")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(data)
}

func CosensawsData(w http.ResponseWriter, r *http.Request) {
    data, err := fetchFileData("Cosensaws")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(data)
}

func FlukeData(w http.ResponseWriter, r *http.Request) {
    data, err := fetchFileData("Fluke")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(data)
}

func KinkelderData(w http.ResponseWriter, r *http.Request) {
    data, err := fetchFileData("Kinkelder")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(data)
}