package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elue-dev/usersapi/controllers"
	"github.com/elue-dev/usersapi/database"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initializeRouter() {
	err := godotenv.Load()

	if err != nil {
	  log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")

	fmt.Println("Go server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(allowedOrigins)(router)))

}

func main() {
	database.InitialMigration()
	initializeRouter()
}