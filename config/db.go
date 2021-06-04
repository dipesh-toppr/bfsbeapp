package config

import (
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database is the object uses by the models for accessing
// database tables and executing queries.
var Database *gorm.DB

type Slot struct {
	ID            uint
	TeacherId     uint
	AvailableSlot uint
	IsBooked      uint
}

type Booked struct {
	ID        uint
	StudentId uint
	SlotId    uint
}

func init() {
	var err error
	Database, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=my_db password=9110SOMU@ sslmode=disable")

	if err != nil {
		panic(err)
	}
	//create table named slots
	if !Database.HasTable(&Slot{}) {
		Database.CreateTable(&Slot{})
	}
	//create table named booked
	if !Database.HasTable(&Booked{}) {
		Database.CreateTable(&Booked{})
	}
	// set this to 'true' to see sql logs
	Database.LogMode(true)

	fmt.Println("Database connection successful.")
}
