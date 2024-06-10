package routers

import (
    "github.com/gorilla/mux"
    "back/controllers"
)

func ScrapperRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/mitutoyo", controllers.Mitutoyo).Methods("GET")
    router.HandleFunc("/cosensaws", controllers.Cosensaws).Methods("GET")
    router.HandleFunc("/fluke", controllers.Fluke).Methods("GET")
    router.HandleFunc("/kinkelder", controllers.Kinkelder).Methods("GET")

    return router
}