package managers

import (
	"errors"

	//"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dipesh-toppr/bfsbeapp/models"
)

///function to validate the timing
func ValidateTime(w http.ResponseWriter, r *http.Request) (uint, error) {
	tim := r.URL.Query()["time"][0]
	hr := time.Now().Hour()
	mn := time.Now().Minute()
	pt, err := strconv.Atoi(tim)
	if err != nil {
		http.Error(w, "invalid time", http.StatusBadRequest)
		return uint(pt), errors.New("invalid time")
	}
	if pt <= hr || (pt == hr+1 && mn > 0) {
		http.Error(w, "booking not allowed at this time", http.StatusBadRequest)
		return uint(pt), errors.New("booking not allowed at this time")
	}
	return uint(pt), nil
}

//check for available teacher
func AvailSlot(tim uint) (uint, error) {
	var slot models.Slot
	tmp := Database.Where("available_slot = ? AND is_booked = ?", tim, false).First(&slot)

	if tmp.Error != nil {
		return slot.ID, errors.New("no availbale slot at this time")
	}
	return slot.ID, nil
}

//book the slot
func BookSlot(stid uint, slid uint) (uint, error) {
	var booked models.Booked
	booked.StudentId = stid
	booked.SlotId = slid
	result1 := Database.Model(&models.Slot{}).Where("id = ? ", slid).Update("is_booked", true)
	if result1.Error != nil {
		return 0, errors.New(result1.Error.Error())
	}
	result2 := Database.Create(&booked)
	if result2.Error != nil {
		return 0, errors.New(result2.Error.Error())
	}
	return booked.ID, nil
}

//read bookings
func ReadBooked(r *http.Request) (models.Slot, bool) {
	var slot models.Slot
	bid := r.URL.Query()["bid"][0] //get the booking id
	bookingId, err := strconv.Atoi(bid)
	if err != nil {
		return slot, false
	}
	var booked models.Booked
	result := Database.Where("id = ?", uint(bookingId)).Find(&booked)
	if result.Error != nil {
		return slot, false
	}
	slotid := booked.SlotId
	result1 := Database.Where("id = ?", slotid).Find(&slot)
	if result1.Error != nil {
		return slot, false
	}
	return slot, true
}

//check for already booked slot at a given time
func IsAlreadyBooked(uid uint, tim uint) bool {
	var booked []models.Booked
	var slotid []uint
	var slot models.Slot
	Database.Where("student_id = ?", uid).Find(&booked)
	for _, val := range booked {
		slotid = append(slotid, val.SlotId)
	}
	result := Database.Where("available_slot = ?", tim).Find(&slot, slotid)
	return result.Error == nil
}

// find booking
func FindBookingAndDelete(bid uint, sid uint) error {
	var booked models.Booked
	result1 := Database.Where("id = ? AND student_id= ?", bid, sid).Find(&booked)
	if result1.Error != nil {
		return errors.New("invalid booking id")
	}
	result2 := Database.Where("id = ?", booked.ID).Delete(&booked)
	if result2.Error != nil {
		return errors.New("internal database error")
	}
	result3 := Database.Model(&models.Slot{}).Where("id = ? ", booked.SlotId).Update("is_booked", false)
	if result3.Error != nil {
		return errors.New("internal database error")
	}
	return nil
}
