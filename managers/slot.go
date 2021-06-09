package managers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/models"
)

func SaveSlot(w http.ResponseWriter, r *http.Request, id int) (models.Slot, error) {
	//to validate the data
	s, err := validateSlotForm(r, uint(id))
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

func validateSlotForm(r *http.Request, id uint) (models.Slot, error) {

	as := r.FormValue("available_slot")
	s := models.Slot{TeacherId: id}

	//validating the slot timing it should be between 0 and 24
	a, err := strconv.Atoi(as)
	if err != nil || a > 24 || a < 1 {
		return models.Slot{}, errors.New("available_slot should be a number between 1 and 23")
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

func FindSlotWithInfo(w http.ResponseWriter, teachID, newSlot int) (s models.Slot, e error) {
	e = Database.Find(&models.Slot{}, "teacher_id=? AND available_slot=?", teachID, newSlot).Error
	if (s == models.Slot{}) {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}
func UpdateSlot(w http.ResponseWriter, slotId string, newSlot int) (e error) {
	e = Database.Model(&models.Slot{}).Where("id=?", slotId).Update("available_slot", newSlot).Error
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}
func DistinctSlots(w http.ResponseWriter) (slots []models.Slot, e error) {
	e = Database.Raw("SELECT * FROM slots WHERE is_booked=? ORDER BY available_slot", 0).Scan(&slots).Error
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}
func DeleteSlot(w http.ResponseWriter, s models.Slot) (e error) {
	e = Database.Delete(&s).Error
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	return
}
