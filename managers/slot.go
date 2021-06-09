package managers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/models"
)

func SaveSlot(r *http.Request, id int) (models.Slot, error) {
	//to validate the data
	s, err := validateSlotForm(r, uint(id))
	if err != nil {
		return s, err
	}
	s.TeacherId = uint(id)
	if Database.Create(&s).Error != nil {
		return s, errors.New("unable to process the transaction")
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
		return models.Slot{}, errors.New("models.Slot already exits")
	}
	return s, nil
}
