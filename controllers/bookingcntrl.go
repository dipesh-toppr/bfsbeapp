package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/managers"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

// search teahcer for specific timing

func SearchTeacher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)
		if e != nil {
			return
		}
		utype := managers.UserType(uint(id)) //checking type of user
		if utype != 1 {
			http.Error(w, "you are not allowed to book session!", http.StatusBadRequest)
			return
		}
		tim, err := managers.ValidateTime(w, r)
		if err != nil {
			return
		}
		//check for already booking at this time
		if managers.IsAlreadyBooked(uint(id), tim) {
			http.Error(w, "you have already booked a session at this time", http.StatusBadRequest)
			return
		}
		//check for available slot at time tim
		slot, err := managers.AvailSlot(tim)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//book the slot
		bookid, err := managers.BookSlot(uint(id), uint(slot))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			return
		}
		bid := r.URL.Query()["bid"][0]
		bkid, err2 := strconv.Atoi(bid)
		if err2 != nil {
			http.Error(w, "student_id or booking_id should be a number", http.StatusBadRequest)
			return
		}
		ok := managers.FindBookingAndDelete(uint(bkid), uint(id))
		if ok != nil {
			http.Error(w, ok.Error(), http.StatusBadRequest)
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
			return
		}
		slot, ok := managers.ReadBooked(r)
		if !ok {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(slot)
		w.WriteHeader(http.StatusOK)
	}
}
