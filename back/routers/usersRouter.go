package routers

import (
    "github.com/gorilla/mux"

    "back/controllers"
)

func UserRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
    router.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
    router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
    router.HandleFunc("/users", controllers.DeleteAllUsers).Methods("DELETE")

    return router
}