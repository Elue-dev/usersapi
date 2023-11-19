package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteuser).Methods("DELETE")

	fmt.Println("Go server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func main() {
	InitialMigration()
	initializeRouter()
}