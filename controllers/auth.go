package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
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

	 createdUser := models.User{}
	 result.Scan(&createdUser)
	 token, err := helpers.GenerateToken(strconv.Itoa(int(createdUser.ID)))
	 if err != nil {
		log.Fatalf("Error generating: %v\n", result.Error)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":  helpers.DatabaseUserToUserModel(user),
		"token": token,
	})
}

func SignUp(w http.ResponseWriter, r *http.Request) {
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
	
	 json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":  helpers.DatabaseUserToUserModel(user),

	})
}