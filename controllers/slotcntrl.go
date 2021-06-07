package controllers

import (
	"encoding/json"
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
			http.Error(w, "authentication failed", http.StatusBadRequest)
			return
		}
		teach := models.User{}
		if config.Database.Where("id=?", id).First(&teach).Error != nil {
			http.Error(w, "unable to process the transaction", http.StatusBadGateway)
			return
		}
		if teach.Identity != strconv.Itoa(0) {
			http.Error(w, "Only teacher can add time slots", http.StatusBadGateway)
			return
		}
		//saving the slot in the database
		s, err := models.SaveSlot(r, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("slot created sucessfully\n"))
		json.NewEncoder(w).Encode(s)
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
			http.Error(w, "authentication failed", http.StatusBadRequest)
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
		s.AvailableSlot = uint(newSlot)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s)
	}
}
func GetUniqueSlots(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		//user authentication
		_, e := token.Parsetoken(w, r)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
		}

		slots := []models.Slot{}
		as := make(map[int]bool)

		if e := config.Database.Raw("SELECT * FROM slots WHERE is_booked=? ORDER BY available_slot", 0).Scan(&slots).Error; e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		keys := []int{}
		for _, i := range slots {
			_, ok := as[int(i.AvailableSlot)]
			if !ok {
				as[int(i.AvailableSlot)] = true
				keys = append(keys, int(i.AvailableSlot))
			}
		}

		w.WriteHeader(http.StatusOK)
		//writing the json response to the response writter
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keys)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
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
		json.NewEncoder(w).Encode(s)
		w.WriteHeader(http.StatusAccepted)
	}
}
