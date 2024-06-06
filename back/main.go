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
    mainRouter.Handle("/users/", middleware.LoggingMiddleware(middleware.ValidateUserInput(userRouter)))
    mainRouter.Handle("/scrapper/", middleware.LoggingMiddleware(http.StripPrefix("/scrapper", scrapperRouter)))
    mainRouter.Handle("/data/", middleware.LoggingMiddleware(dataRouter))
    
    corsRouter := middleware.CorsMiddleware(mainRouter)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    log.Printf("Listening on port %s...\n", port)
    log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}