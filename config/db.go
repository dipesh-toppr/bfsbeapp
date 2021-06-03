package config

import (
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database is the object uses by the models for accessing
// database tables and executing queries.
var Database *gorm.DB

func init() {
	var err error
	Database, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=bfsbedata password=daddycool sslmode=disable")

	if err != nil {
		panic(err)
	}

	// set this to 'true' to see sql logs
	Database.LogMode(true)

	fmt.Println("Database connection successful.")
}
