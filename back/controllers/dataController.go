package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"back/models"
	"back/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func fetchFileData(company string) ([]models.ProductData, error) {
    var results []models.ProductData
    filter := bson.M{"company": company}
    cursor, err := utils.FileDataCollection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var fileData models.ProductData
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

func RemoveAllMitutoyoData(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{"companyName": "Mitutoyo"}

	// DeleteMany ejecuta la eliminación de todos los documentos que coincidan con el filtro
	result, err := utils.FileDataCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Informar sobre la cantidad de documentos eliminados
	response := map[string]interface{}{
		"message":          "Deleted all Mitutoyo data successfully",
		"deletedDocuments": result.DeletedCount,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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