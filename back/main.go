package main

import (
    "log"
    "net/http"
    "os"

    "back/middleware"
    "back/routers"
)

func main() {
    userRouter := routers.UserRouter()
    scrapperRouter := routers.ScrapperRouter()
    dataRouter := routers.ShowData()

    mainRouter := http.NewServeMux()
    // Working tested
    mainRouter.Handle("/users/", middleware.LoggingMiddleware(middleware.ValidateUserInput(http.StripPrefix("/users", userRouter))))
    // Working, tested
    mainRouter.Handle("/scrapper/", middleware.LoggingMiddleware(http.StripPrefix("/scrapper", scrapperRouter)))
    // Working, tested
    mainRouter.Handle("/data/", middleware.LoggingMiddleware(http.StripPrefix("/data", dataRouter)))
    
    corsRouter := middleware.CorsMiddleware(mainRouter)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    log.Printf("Listening on port %s...\n", port)
    log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}