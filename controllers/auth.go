package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/elue-dev/usersapi/database"
	"github.com/elue-dev/usersapi/helpers"
	"github.com/elue-dev/usersapi/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
    var payload models.Login
	json.NewDecoder(r.Body).Decode(&payload)

    result := database.DB.Where("email = ?", payload.Email).First(&user)

    if result.Error != nil {
        log.Printf("User not found for email %s", payload.Email)
        http.Error(w, "Invalid credentials provided", http.StatusUnauthorized)
        return
    }

	 passwordIsCorrect := helpers.ComparePasswordWithHash(user.Password, payload.Password)
	 if !passwordIsCorrect {
		 log.Printf("Invalid password for email %s", payload.Email)
		 http.Error(w, "Invalid credentials provided", http.StatusUnauthorized)
		 return
	 }

	json.NewEncoder(w).Encode(helpers.DatabaseUserToUserModel(user))
	
}