package main

import (
    "github.com/gorilla/mux"
    "back/controllers"
)

func RegisterRoutes() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
    router.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
    router.HandleFunc("/generate-csv", controllers.GenerateCSV).Methods("GET")

    return router
}