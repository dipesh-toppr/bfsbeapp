package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/config"
	"github.com/dipesh-toppr/bfsbeapp/models"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

func AddSlot(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		//user authentication
		id, e := token.Parsetoken(w, r)
		if e != nil {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}
		//saving the slot in the database
		s, err := models.SaveSlot(r, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("slot created sucessfully\n", s)
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func GetUserSlots(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		//user authentication
		id, e := token.Parsetoken(w, r)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
		}

		slots := []models.Slot{}
		//getting the slots of the user
		if e = config.Database.Find(&slots, "teacher_id=?", id).Error; e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

		//writing the json response to the response writter
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(slots)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func UpdateSlot(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		//user authentication
		id, e := token.Parsetoken(w, r)
		slotId := r.FormValue("slot_id")
		newSlot, _ := strconv.Atoi(r.FormValue("new_slot"))
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
		}
		s := models.Slot{}
		config.Database.Find(&s, "id=?", slotId)
		teachID := s.TeacherId
		if teachID != uint(id) {
			http.Error(w, "authentication failed", http.StatusBadRequest)
			return
		}
		if config.Database.Find(&models.Slot{}, "teacher_id=? AND available_slot=?", teachID, newSlot).Error == nil {
			http.Error(w, "Slot already exists", http.StatusBadRequest)
			return
		}
		if e := config.Database.Model(&models.Slot{}).Where("id=?", slotId).Update("available_slot", newSlot).Error; e != nil {
			http.Error(w, e.Error(), http.StatusExpectationFailed)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		id, e := token.Parsetoken(w, r)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
		}
		slotId := r.FormValue("DEL_slot")
		s := models.Slot{}
		config.Database.Find(&s, "id=?", slotId)
		teachId := s.TeacherId
		if teachId != uint(id) {
			http.Error(w, "authentication failed", http.StatusBadRequest)
			return
		}
		if e := config.Database.Delete(&s).Error; e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
