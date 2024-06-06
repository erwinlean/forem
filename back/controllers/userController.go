package controllers

import (
    "log"
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt" 
    "go.mongodb.org/mongo-driver/bson"
    "back/models"
    "back/utils"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("forem") 

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var result models.User
    err = utils.UserCollection.FindOne(context.TODO(), bson.M{
        "$or": []bson.M{
            {"username": user.Username},
            {"email": user.Email},
        },
    }).Decode(&result)
    if err != nil {
        http.Error(w, "Invalid username, email, or password", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
    if err != nil {
        http.Error(w, "Invalid username, email, or password", http.StatusUnauthorized)
        return
    }

    result.LoginDates = append(result.LoginDates, time.Now())

    _, err = utils.UserCollection.ReplaceOne(context.TODO(), bson.M{"_id": result.ID}, result)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // jwt
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": result.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    log.Println("token generated on the loggin...")
    log.Println(tokenString)
    log.Println(jwtKey)

    w.Header().Set("Authorization", "Bearer "+ tokenString)

    json.NewEncoder(w).Encode(struct {
        User  models.User `json:"user"`
        Token string      `json:"token"`
    }{User: result, Token: tokenString})
}

// Just for test and delete the data at the DB
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var existingUser models.User
    err = utils.UserCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
    if err == nil {
        http.Error(w, "Email already registered", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    user.CreatedAt = time.Now()

    _, err = utils.UserCollection.InsertOne(context.TODO(), user)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func DeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	result, err := utils.UserCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Println("Error deleting all users:", err)
		
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Deleted %d users\n", result.DeletedCount)
	
	w.WriteHeader(http.StatusNoContent)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	cursor, err := utils.UserCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("Error retrieving users:", err)
		
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		log.Println("No users found")
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding users to JSON:", err)
		
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}