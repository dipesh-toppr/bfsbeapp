package managers

import (
	"fmt"

	"github.com/dipesh-toppr/bfsbeapp/models"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database is the object uses by the managers for accessing
// database tables and executing queries.
var Database *gorm.DB

func init() {
	var err error
	Database, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=6677 sslmode=disable")

	if err != nil {
		panic(err)
	}
	if !Database.HasTable(&models.User{}) {
		Database.CreateTable(&models.User{})
	}
	//create table named slots
	if !Database.HasTable(&models.Slot{}) {
		Database.CreateTable(&models.Slot{})
	}
	//create table named booked
	if !Database.HasTable(&models.Booked{}) {
		Database.CreateTable(&models.Booked{})
	}
	// set this to 'true' to see sql logs
	Database.LogMode(true)

	fmt.Println("Database connection successful.")
}
