package middleware

import (
    "log"
    "net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: Method=%s, URL=%s", r.Method, r.URL.String())
        next.ServeHTTP(w, r)
    })
}