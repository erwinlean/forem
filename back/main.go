package main

import (
    "log"
    "net/http"
)

func main() {
    router := RegisterRoutes()
    
    log.Println("server on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}