package models

import (
	"errors"
	"fmt"
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

//read bookings
func ReadBooked(r *http.Request) (config.Slot, bool) {
	var slot config.Slot
	bid := r.URL.Query()["bid"][0] //get the booking id
	bookingId, err := strconv.Atoi(bid)
	if err != nil {
		return slot, false
	}
	var booked config.Booked
	result := config.Database.Where("id = ?", uint(bookingId)).Find(&booked)
	if result.Error != nil {
		return slot, false
	}
	slotid := booked.SlotId
	result1 := config.Database.Where("id = ?", slotid).Find(&slot)
	if result1.Error != nil {
		return slot, false
	}
	return slot, true
}

//check for already booked slot at a given time
func IsAlreadyBooked(uid uint, tim uint) bool {
	var booked []config.Booked
	config.Database.Where("student_id = ?", uid).Find(&booked)
	for _, val := range booked {
		var slot config.Slot
		config.Database.Where("id = ?", val.SlotId).Find(&slot)
		if slot.AvailableSlot == tim {
			return true
		}
	}
	return false
}
func ReadStudents(r *http.Request) ([]config.User, bool) {

	var stud []config.User
	result := config.Database.Where("identity = ?", uint(1)).Find(&stud)
	if result.Error != nil {
		return stud, false
	}
	fmt.Print(stud)
	return stud, true
}

func ReadTeachers(r *http.Request) ([]config.User, bool) {

	var teach []config.User
	result := config.Database.Where("identity = ?", uint(0)).Find(&teach)
	if result.Error != nil {
		return teach, false
	}
	fmt.Print(teach)
	return teach, true
}

func ReadAdminBooked(r *http.Request) ([]config.Booked, bool) {

	var booked []config.Booked
	result := config.Database.Find(&booked)
	if result.Error != nil {
		return booked, false
	}

	return booked, true
}
