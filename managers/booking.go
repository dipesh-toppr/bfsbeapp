package managers

import (
	"errors"
	"fmt"

	//"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dipesh-toppr/bfsbeapp/models"
)

///function to validate the timing
func ValidateTime(w http.ResponseWriter, r *http.Request) (string, error) {
	tim := r.URL.Query()["date"][0]
	_, err := time.Parse("2006-Jan-02", tim)
	if err != nil {
		fmt.Println(err)
		return tim, errors.New("invalid date")
	}
	return tim, nil
}

//check for available teacher
func AvailSlot(date string) ([]models.Slot, error) {
	var slots []models.Slot
	CurrTime := time.Now().Hour()
	minute := time.Now().Minute()
	var TimeToSearch uint
	TimeToSearch = uint(CurrTime + 2)
	if minute == 0 {
		TimeToSearch = uint(CurrTime + 1)
	}
	tmp := Database.Where("available_slot > ? AND is_booked = ? AND date = ?", TimeToSearch, false, date).Find(&slots)

	if tmp.Error != nil {
		return slots, errors.New("no availbale slot at this time")
	}
	return slots, nil
}

//book the slot
func BookSlot(stid uint, slid uint) (uint, error) {
	var booked models.Booked
	booked.StudentId = stid
	booked.SlotId = slid
	var slot models.Slot
	result1 := Database.Where("id = ? AND is_booked = ?", slid, false).Find(&slot)
	if result1.RowsAffected == 0 {
		return 0, errors.New("slot id not found or this slot is already booked")
	}
	result2 := Database.Model(&models.Slot{}).Where("id = ? ", slid).Update("is_booked", true)
	if result2.Error != nil {
		return 0, errors.New(result2.Error.Error())
	}
	result3 := Database.Create(&booked)
	if result3.Error != nil {
		return 0, errors.New(result3.Error.Error())
	}
	return booked.ID, nil
}

//read bookings
func ReadBooked(sid uint, r *http.Request) (models.Slot, bool) {
	var slot models.Slot
	bid := r.URL.Query()["bid"][0] //get the booking id
	bookingId, err := strconv.Atoi(bid)
	if err != nil {
		return slot, false
	}
	var booked models.Booked
	result := Database.Where("id = ? AND student_id = ?", uint(bookingId), sid).Find(&booked)
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
