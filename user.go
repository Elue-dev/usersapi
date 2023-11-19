package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error
const dsn = "host=localhost user=postgres password=pasport dbname=usersapi port=5432 sslmode=disable"

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type CustomUser struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `json:"deleted_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}


func InitialMigration() {
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to DB")
	}

	fmt.Println("Connected to Postgres DB")

	DB.AutoMigrate(&User{})
}


func getUsers(w http.ResponseWriter, r *http.Request) {

}

func getUser(w http.ResponseWriter, r *http.Request) {

}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	 hashedPassword, err := hashPassword(user.Password)
	 if err != nil {
		fmt.Println("Could not hash user password", err)
		 return
	 }

	 user.Password = hashedPassword

	DB.Create(&user)

	customUser := CustomUser{
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

func updateUser(w http.ResponseWriter, r *http.Request) {

}

func deleteuser(w http.ResponseWriter, r *http.Request) {

}

