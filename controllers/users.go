package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/elue-dev/usersapi/database"
	"github.com/elue-dev/usersapi/helpers"
	"github.com/elue-dev/usersapi/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)




var err error


func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []models.User
	database.DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var user models.User
	result := database.DB.First(&user, params["id"])

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "User with id " +  params["id"] + " not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	 hashedPassword, err := helpers.HashPassword(user.Password)
	 if err != nil {
		fmt.Println("Could not hash user password", err)
		 return
	 }

	 user.Password = hashedPassword
	

	 database.DB.Create(&user)

	customUser := models.CustomUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
	
	json.NewEncoder(w).Encode(customUser)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var user models.User
	result := database.DB.First(&user, params["id"])

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "User with id " +  params["id"] + " not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if _, ok := requestBody["password"]; ok {
		http.Error(w, "Updating password directly is not allowed", http.StatusBadRequest)
		return
	}

	json.NewDecoder(r.Body).Decode(&user)
	database.DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var user models.User
	result := database.DB.First(&user, params["id"])

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "User with id " +  params["id"] + " not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	database.DB.Delete(&user, params["id"])
	json.NewEncoder(w).Encode("User has been deleted")

}

