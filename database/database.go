package database

import (
	"fmt"

	"github.com/elue-dev/usersapi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error
const dsn = "host=localhost user=postgres password=pasport dbname=usersapi port=5432 sslmode=disable"


func InitialMigration() {
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Could not connect to DB")
	}

	fmt.Println("Connected to Postgres DB")

	DB.AutoMigrate(&models.User{})
}

