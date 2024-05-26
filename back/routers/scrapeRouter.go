package routers

import (
    "github.com/gorilla/mux"

    "back/controllers"
)

func ScrapperRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/mitutoyoScrapper", controllers.Mitutoyo).Methods("GET")
    router.HandleFunc("/cosensawsScrapper", controllers.Cosensaws).Methods("GET")
    router.HandleFunc("/flukeScrapper", controllers.Fluke).Methods("GET")
    router.HandleFunc("/kinkelderScrapper", controllers.Kinkelder).Methods("GET")

    return router
}