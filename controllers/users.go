package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/elue-dev/usersapi/database"
	"github.com/elue-dev/usersapi/helpers"
	"github.com/elue-dev/usersapi/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)


func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []models.User
	database.DB.Find(&users)
	// json.NewEncoder(w).Encode(helpers.DatabaseUsersArrToUserModel(users))
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"total": len(users),
		"users":  helpers.DatabaseUsersArrToUserModel(users),
	})
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":  helpers.DatabaseUserToUserModel(user),
	})
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

	//TODO: Update user db query
	
	json.NewDecoder(r.Body).Decode(&user)
	database.DB.Save(&user)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":  helpers.DatabaseUserToUserModel(user),
	})
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

	userFromToken, err := helpers.GetUserFromToken(r)
	if err != nil {
		http.Error(w, "Error getting token", http.StatusInternalServerError)
		return
	}

	if strconv.FormatUint(uint64(userFromToken.ID), 10) != params["id"] {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "You can only delete your account",
		})
		return
	}

	database.DB.Delete(&user, params["id"])
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "User has been deleted",
	})
}




