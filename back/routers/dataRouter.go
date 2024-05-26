package routers

import (
    "github.com/gorilla/mux"

    "back/controllers"
)

func ShowData() *mux.Router {
    router := mux.NewRouter()

    // Data getter
    router.HandleFunc("/mitutoyoData", controllers.MitutoyoData).Methods("GET")
    router.HandleFunc("/cosensawsData", controllers.CosensawsData).Methods("GET")
    router.HandleFunc("/flukeData", controllers.FlukeData).Methods("GET")
    router.HandleFunc("/kinkelderData", controllers.KinkelderData).Methods("GET")

    return router
}