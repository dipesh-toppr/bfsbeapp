package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dipesh-toppr/bfsbeapp/config"
	"github.com/dipesh-toppr/bfsbeapp/models"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

func AddSlot(w http.ResponseWriter, r *http.Request) {

	//user authentication
	e, id := token.Parsetoken(w, r)
	if e != nil {
		http.Error(w, "unauthorized request", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodPost {
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

	//user authentication
	e, id := token.Parsetoken(w, r)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}

	slots := []models.Slot{}
	if r.Method == http.MethodGet {
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
