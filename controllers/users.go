package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"context"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
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
	json.NewEncoder(w).Encode(helpers.DatabaseUsersArrToUserModel(users))
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

	json.NewEncoder(w).Encode(helpers.DatabaseUserToUserModel(user))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        log.Fatalf("Failed to parse form data: %v", err)
        return
    }

	// json.NewDecoder(r.Body).Decode(&user)
	user.FirstName = r.FormValue("first_name")
    user.LastName = r.FormValue("last_name")
    user.Email = r.FormValue("email")
    user.Password = r.FormValue("password")

	hashedPassword, err := helpers.HashPassword(user.Password)
	 if err != nil {
		fmt.Println("Could not hash user password", err)
		 return
	 }

	 user.Password = hashedPassword

	 file, _, err := r.FormFile("avatar")
	 if err != nil {
		 log.Fatalf("Failed to get avatar from form: %v", err)
		 return
	 }
	 defer file.Close()
 
	 cld, err := cloudinary.New()
	 if err != nil {
		 log.Fatalf("Failed to initialize Cloudinary: %v", err)
		 return
	 }
 
	 var ctx = context.Background()
	 uploadResult, err := cld.Upload.Upload(
		 ctx,
		 file,
         uploader.UploadParams{PublicID: "avatar"})
 
	 if err != nil {
		 log.Fatalf("Failed to upload file: %v\n", err)
		 return
	 }
  
	 user.Avatar = uploadResult.SecureURL

	 database.DB.Create(&user)
	
	json.NewEncoder(w).Encode(helpers.DatabaseUserToUserModel(user))
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
	json.NewEncoder(w).Encode(helpers.DatabaseUserToUserModel(user))
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




