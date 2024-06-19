package routers

import (
    "github.com/gorilla/mux"
    "back/controllers"
)

func ShowData() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/mitutoyo", controllers.MitutoyoData).Methods("GET")
    router.HandleFunc("/cosensaws", controllers.CosensawsData).Methods("GET")
    router.HandleFunc("/fluke", controllers.FlukeData).Methods("GET")
    router.HandleFunc("/kinkelder", controllers.KinkelderData).Methods("GET")

    router.HandleFunc("/mitutoyo", controllers.RemoveAllMitutoyoData).Methods("DELETE")

    return router
}