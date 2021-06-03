package models

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/config"
)

type Slot struct {
	ID            uint
	TeacherId     uint
	AvailableSlot uint
}

func SaveSlot(r *http.Request) (Slot, error) {
	s, err := validateSlotForm(r)
	if err != nil {
		return s, err
	}
	if config.Database.Create(&s).Error != nil {
		return s, errors.New("unable to process the transaction")
	}
	return s, nil
}

func validateSlotForm(r *http.Request) (Slot, error) {
	t := r.FormValue("teacher_id")
	as := r.FormValue("available_slot")
	s := Slot{}
	ti, err := strconv.Atoi(t)
	if err != nil {
		return s, errors.New("invalid teacher_id")
	}
	a, err := strconv.Atoi(as)
	if err != nil || a > 24 || a < 0 {
		return Slot{}, errors.New("available_slot should be a number between 0 and 24")
	}
	s.TeacherId = uint(ti)
	s.AvailableSlot = uint(a)
	return s, nil
}
