package routers

import (
    "log"
    
    "github.com/gorilla/mux"
    "back/controllers"
)

func ScrapperRouter() *mux.Router {
    router := mux.NewRouter()

    log.Println("Setting up scrapper routes")

    router.HandleFunc("/mitutoyo", controllers.Mitutoyo).Methods("GET")
    router.HandleFunc("/cosensaws", controllers.Cosensaws).Methods("GET")
    router.HandleFunc("/fluke", controllers.Fluke).Methods("GET")
    router.HandleFunc("/kinkelder", controllers.Kinkelder).Methods("GET")

    return router
}