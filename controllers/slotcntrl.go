package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/managers"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

func AddSlot(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		//user authentication
		id, e := token.Parsetoken(w, r)
		if e != nil {
			return
		}
		user, e := managers.FindUserWithId(w, id)
		fmt.Println(user)
		if e != nil {
			return
		}
		if user.Identity != (0) {
			http.Error(w, "Only teacher can add time slots", http.StatusBadRequest)
			return
		}
		//saving the slot in the database
		s, err := managers.SaveSlot(w, r, id)
		if err != nil {
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
			return
		}
		slots, e := managers.FindSlotWithTeacherId(w, id)
		//getting the slots of the user
		if e != nil {
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
		date := r.FormValue("date")
		newSlot, _ := strconv.Atoi(r.FormValue("new_slot"))
		if e != nil {
			return
		}
		s, e := managers.FindSlotWithId(w, slotId)
		if e != nil {
			return
		}
		if s.IsBooked {
			http.Error(w, "you cannot update a booked slot, delete and create new slot", http.StatusForbidden)
			return
		}
		teachID := s.TeacherId
		if teachID != uint(id) {
			http.Error(w, "you can only update your slots", http.StatusBadRequest)
			return
		}
		_, err := managers.FindSlotWithInfo(w, int(teachID), newSlot, date)
		if err == nil {
			http.Error(w, "Slot already exists", http.StatusBadRequest)
			return
		}
		er := managers.UpdateSlot(w, slotId, date, newSlot)
		if er != nil {
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
			return
		}

		slots, e := managers.DistinctSlots(w)
		as := make(map[int]bool)

		if e != nil {
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
			return
		}
		slotId := r.FormValue("DEL_slot")
		s, e := managers.FindSlotWithId(w, slotId)
		if e != nil {
			return
		}

		teachId := s.TeacherId
		if teachId != uint(id) {
			http.Error(w, "you cannot delete other teacher's slots", http.StatusForbidden)
			return
		}
		er := managers.DeleteSlot(w, s)
		if er != nil {
			return
		}
		json.NewEncoder(w).Encode(s)
		w.WriteHeader(http.StatusAccepted)
	}
}
