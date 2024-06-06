package middleware

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "github.com/golang-jwt/jwt"
)

var EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func IsValidEmail(email string) bool {
    return EmailRegex.MatchString(email)
}

func ContainsMaliciousContent(input string) bool {
    return strings.Contains(input, "<script>") || strings.Contains(input, "</script>") || strings.Contains(input, "|") || strings.Contains(input, "\n")|| strings.Contains(input, "INSERT")|| strings.Contains(input, "UPDATE")|| strings.Contains(input, "<>")
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

func ValidateToken(tokenString, secretKey string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    return token, nil
}