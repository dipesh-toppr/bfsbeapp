package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dipesh-toppr/bfsbeapp/config"
)

///function to validate the timing
func ValidateTime(r *http.Request) (uint, error) {
	tim := r.URL.Query()["time"][0]
	hr := time.Now().Hour()
	mn := time.Now().Minute()
	pt, err := strconv.Atoi(tim)
	if err != nil {
		return uint(pt), errors.New("invalid time")
	}
	if pt <= hr || (pt == hr+1 && mn <= 59) {
		return uint(pt), errors.New("booking not allowed at this time")
	}
	return uint(pt), nil
}

//check for available teacher
func AvailSlot(tim uint) (uint, error) {
	var slot Slot
	tmp := config.Database.Where("available_slot = ? AND is_booked = ?", tim, 0).First(&slot)
	if tmp.Error != nil {
		return slot.ID, errors.New("no availbale slot at this time")
	}
	return slot.ID, nil
}

//book the slot
func BookSlot(stid uint, slid uint) (uint, error) {
	var booked config.Booked
	booked.StudentId = stid
	booked.SlotId = slid
	result1 := config.Database.Model(&Slot{}).Where("id = ? ", slid).Update("is_booked", 1)
	if result1.Error != nil {
		return 0, errors.New(result1.Error.Error())
	}
	result2 := config.Database.Create(&booked)
	if result2.Error != nil {
		return 0, errors.New(result2.Error.Error())
	}
	return booked.ID, nil
}
