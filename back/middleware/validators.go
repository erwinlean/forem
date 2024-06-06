package middleware

import (
    "net/http"
    "regexp"
    "strings"
)

var EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func IsValidEmail(email string) bool {
    return EmailRegex.MatchString(email)
}

func ContainsMaliciousContent(input string) bool {
    // Aquí se puede agregar lógica adicional para verificar contenido malicioso.
    return strings.Contains(input, "<script>") || strings.Contains(input, "</script>")
}

func ValidateUserInput(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        
        email := r.FormValue("email")
        if email != "" && !IsValidEmail(email) {
            http.Error(w, "Invalid email address", http.StatusBadRequest)
            return
        }
        
        for key, values := range r.Form {
            for _, value := range values {
                if ContainsMaliciousContent(value) {
                    http.Error(w, "Malicious content detected in "+key, http.StatusBadRequest)
                    return
                }
            }
        }

        next.ServeHTTP(w, r)
    })
}