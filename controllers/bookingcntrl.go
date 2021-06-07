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

// search teahcer for specific timing

func SearchTeacher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)
		if e != nil {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}
		utype := models.UserType(uint(id)) //checking type of user
		if utype != "1" {
			http.Error(w, "you are not allowed to book session!", http.StatusBadRequest)
			return
		}
		tim, err := models.ValidateTime(r)
		print(tim, " ", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//check for already booking at this time
		if models.IsAlreadyBooked(uint(id), tim) {
			http.Error(w, "you have already booked a session at this time", http.StatusBadRequest)
			return
		}
		//check for available slot at time tim
		slot, err := models.AvailSlot(tim)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//book the slot
		bookid, err := models.BookSlot(uint(id), uint(slot))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//send booking id to the user
		msg := "booking ID : " + fmt.Sprint(bookid)
		w.Write([]byte(msg))
		w.WriteHeader(http.StatusOK)
		return
	}
}

//delete booked slot
func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)
		if e != nil {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}
		bid := r.URL.Query()["bid"][0]
		bkid, err2 := strconv.Atoi(bid)
		if err2 != nil {
			http.Error(w, "student_id OR booking_id should be a number", http.StatusBadRequest)
			return
		}
		var booked config.Booked
		booked.ID = uint(bkid)
		booked.StudentId = uint(id)
		result3 := config.Database.Where("id = ? AND student_id= ?", booked.ID, booked.StudentId).Find(&booked)
		slot := booked.SlotId
		if result3.Error != nil {
			http.Error(w, "Invalid booking ID", http.StatusBadRequest)
			return
		}
		result1 := config.Database.Where("id = ?", booked.ID).Delete(&booked)
		if result1.Error != nil {
			http.Error(w, result1.Error.Error(), http.StatusInternalServerError)
			return
		}
		result2 := config.Database.Model(&models.Slot{}).Where("id = ? ", slot).Update("is_booked", 0)
		if result2.Error != nil {
			http.Error(w, result2.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Booking Deleted!"))
		w.WriteHeader(http.StatusOK)
		return
	}
}

//read the booking using booking id
func ReadBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)
		if e != nil {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}
		slot, ok := models.ReadBooked(r)
		if !ok {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(slot)
		w.WriteHeader(http.StatusOK)
	}
}
