package main

import (
    "github.com/gorilla/mux"

    "back/controllers"
)

func RegisterRoutes() *mux.Router {
    router := mux.NewRouter()

    // Users
    router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
    router.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")

    // Scrapper functions 
    router.HandleFunc("/mitutoyoScrapper", controllers.Mitutoyo).Methods("GET")
    router.HandleFunc("/cosensawsScrapper", controllers.Cosensaws).Methods("GET")
    router.HandleFunc("/flukeScrapper", controllers.Fluke).Methods("GET")
    router.HandleFunc("/kinkelderScrapper", controllers.Kinkelder).Methods("GET")

    // Data getter
    router.HandleFunc("/mitutoyoData", controllers.MitutoyoData).Methods("GET")
    router.HandleFunc("/cosensawsData", controllers.CosensawsData).Methods("GET")
    router.HandleFunc("/flukeData", controllers.FlukeData).Methods("GET")
    router.HandleFunc("/kinkelderData", controllers.KinkelderData).Methods("GET")

    return router
}