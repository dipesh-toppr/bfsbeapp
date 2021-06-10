package managers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dipesh-toppr/bfsbeapp/models"
)

func SaveSlot(w http.ResponseWriter, r *http.Request, id int) (models.Slot, error) {
	//to validate the data
	s, err := validateSlotForm(w, r, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return s, err
	}
	s.TeacherId = uint(id)
	if e := Database.Create(&s).Error; e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return s, errors.New("x")
	}

	return s, nil
}
func validateTime(w http.ResponseWriter, r *http.Request) (string, error) {
	tim := r.FormValue("date")
	_, err := time.Parse("2006-Jan-02", tim)
	if err != nil {
		fmt.Println(err)
		return tim, errors.New("invalid date")
	}
	return tim, nil
}
func validateSlotForm(w http.ResponseWriter, r *http.Request, id uint) (models.Slot, error) {

	as := r.FormValue("available_slot")
	date := r.FormValue("date")
	s := models.Slot{TeacherId: id, Date: date}
	_, err := validateTime(w, r)
	if err != nil {
		return models.Slot{}, errors.New(err.Error())
	}
	a, err := strconv.Atoi(as)
	if err != nil || a > 23 || a < 1 {
		return models.Slot{}, errors.New("time must be between 1 and 23")
	}

	s.AvailableSlot = uint(a)
	//checking if there is already a slot available in the db to avoid duplicate entries
	if Database.Find(&models.Slot{}, s).Error == nil {
		return models.Slot{}, errors.New("slot already exits")
	}
	return s, nil
}
func FindUserWithId(w http.ResponseWriter, id int) (u models.User, e error) {
	e = Database.Find(&u, "id=?", id).Error
	if (u == models.User{}) {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}
func FindSlotWithTeacherId(w http.ResponseWriter, id int) (slots []models.Slot, e error) {
	e = Database.Find(&slots, "teacher_id=?", id).Error
	if len(slots) == 0 {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}

func FindSlotWithId(w http.ResponseWriter, slotId string) (s models.Slot, e error) {
	e = Database.Find(&s, "id=?", slotId).Error
	if (s == models.Slot{}) {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}

func FindSlotWithInfo(w http.ResponseWriter, teachID, newSlot int, date string) (s models.Slot, e error) {
	e = Database.Find(&models.Slot{}, "teacher_id=? AND available_slot=? AND date=?", teachID, newSlot, date).Error
	if (s == models.Slot{}) {
		return
	}
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}
func UpdateSlot(w http.ResponseWriter, slotId, date string, newSlot int) (e error) {
	if newSlot < 1 || newSlot > 23 {
		http.Error(w, "time must be between 1 and 23", http.StatusBadRequest)
		e = errors.New("_")
		return
	}
	e = Database.Model(&models.Slot{}).Where("id=? AND date=?", slotId, date).Update("available_slot", newSlot).Error
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}
func DistinctSlots(w http.ResponseWriter) (slots []models.Slot, e error) {
	e = Database.Raw("SELECT * FROM slots WHERE is_booked=? ORDER BY available_slot", 0).Scan(&slots).Error
	if len(slots) == 0 {
		http.Error(w, "", http.StatusNotFound)
	}
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}

// func validateTime(pt int) (uint, error) {
// 	hr := time.Now().Hour()
// 	mn := time.Now().Minute()
// 	if pt == hr+1 && mn > 0 {
// 		return uint(pt), errors.New("you can only delete a slot if its one hour ahead current time")
// 	}
// 	return uint(pt), nil
// }
func DeleteSlot(w http.ResponseWriter, s models.Slot) (e error) {
	slots := []models.Slot{}
	if er := Database.Find(&slots, "available_slot=? AND is_booked=?", s.AvailableSlot, 0).Error; er != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	if len(slots) == 0 {
		http.Error(w, "you cannot delete this slot as there is no other teacher to cover for you", http.StatusForbidden)
		e = errors.New("_")
		return
	}
	// tim, err := validateTime(int(s.AvailableSlot))
	// print(tim, " ", err)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	e = errors.New("_")
	// 	return
	// }
	slot := slots[0]
	fmt.Println(slot)
	if Database.Model(&models.Slot{}).Where("id = ? ", slot.ID).Update("is_booked", true).Error != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		e = errors.New("_")
		return
	}
	if Database.Model(&models.Booked{}).Where("slot_id=?", s.ID).Update("slot_id", slot.ID).Error != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		e = errors.New("_")
		return
	}
	e = Database.Delete(&s).Error
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		e = errors.New("_")
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
