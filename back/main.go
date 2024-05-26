package main

import (
    "log"
    "net/http"

    "back/routers"
)

func main() {
    userRouter := routers.UserRouter()
    scrapperRouter := routers.ScrapperRouter()
    dataRouter := routers.ShowData()

    // Show delete the warning/error for not use the routers
    log.Print(scrapperRouter, dataRouter)
    
    log.Println("Port :8080")
    log.Fatal(http.ListenAndServe(":8080", userRouter))
    log.Fatal(http.ListenAndServe(":8080", scrapperRouter))
    log.Fatal(http.ListenAndServe(":8080", dataRouter))
}