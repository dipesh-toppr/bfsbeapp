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

func Admin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		idtodisable := r.FormValue("idtodisable") //obtaining id to disable as input

		id, e := token.Parsetoken(w, r) //finding active user mail and if he is logged in

		// print(e, "  ", mail) //for debugging

		if e != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		user, ok := models.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusForbidden)
			fmt.Println("no user found")
			return
		}

		uid := user.Identity //uid is identity of user ie stud,tech,admin,superadmin

		fmt.Println("This is uid ", uid)

		//finding the user detials to check his/her role

		utodisable, ok := models.FindUserFromId(idtodisable)

		if !ok {
			http.Error(w, "no user found", http.StatusForbidden)
			fmt.Println("no user found")
			return
		}

		if utodisable.Identity < "2" { ///means he is stud or teacher so can be made inactive my both admin and super admin

			if uid >= "2" {
				//if user iddentity is>= 2  means that active user is an admin or super admin & has rights to make any user inactive
				u := models.MakeInactive(idtodisable)
				fmt.Print(u)
				w.Write([]byte("user disabled\n"))
				json.NewEncoder(w).Encode(u)
				return

			} else {
				http.Error(w, "You do not have the rights to make admin inactive", http.StatusBadRequest)
				fmt.Print("You do not have the rights to make user inactive")
			}
		} else if utodisable.Identity == "2" { //request to disable admin do only super admin can do so

			if uid == "3" { //identity  of superadmin  kept 3
				u := models.MakeInactive(idtodisable)
				fmt.Print(u)
				w.Write([]byte("user disabled\n"))
				json.NewEncoder(w).Encode(u)
				return
			} else {
				http.Error(w, "You do not have the rights to make admin inactive", http.StatusBadRequest)
			}
		}
		if utodisable.Identity == "3" {
			http.Error(w, "you cannot disable super admin", http.StatusBadRequest)
		}
		// http.Redirect(w, r, "/", http.StatusOK)
	}

}

//admin read all bookings

func ReadAllBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)

		user, ok := models.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusForbidden)
			fmt.Println("no user found")
			return
		}

		//uid := user.Identity
		fmt.Println(user)
		if e != nil || user.Identity < "2" {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}

		slot, ok := models.ReadAdminBooked(r)
		if !ok {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}
		if len(slot) > 0 {
			json.NewEncoder(w).Encode(slot)

		} else {

			w.Write([]byte("No bookings found"))

		}

		w.WriteHeader(http.StatusOK)
	}
}

func ReadAllTeachers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)

		user, ok := models.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusForbidden)
			fmt.Println("no user found")
			return
		}

		//uid := user.Identity
		fmt.Println(user)

		if e != nil || user.Identity < "2" {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}

		teachers, ok := models.ReadTeachers(r)
		if !ok {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}

		if len(teachers) > 0 {
			json.NewEncoder(w).Encode(teachers)

		} else {

			w.Write([]byte("No teachers found"))

		}

		w.WriteHeader(http.StatusOK)

	}
}

func ReadAllStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)

		user, ok := models.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusForbidden)
			fmt.Println("no user found")
			return
		}

		//uid := user.Identity
		fmt.Println(user)

		if e != nil || user.Identity < "2" {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}

		students, ok := models.ReadStudents(r)
		if !ok {
			http.Error(w, "not found", http.StatusBadRequest)
			return
		}

		if len(students) > 0 {
			json.NewEncoder(w).Encode(students)

		} else {

			w.Write([]byte("No Students found"))

		}

		w.WriteHeader(http.StatusOK)
	}
}

//admin delete booking

func AdminDeleteBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		id, e := token.Parsetoken(w, r) //finding active user mail and if he is logged in

		// print(e, "  ", mail) //for debugging

		if e != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		user, ok := models.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusForbidden)
			fmt.Println("no user found")
			return
		}

		uid := user.Identity

		if uid > "1" {
			bid := r.URL.Query()["bid"][0]
			bkid, err2 := strconv.Atoi(bid)
			if err2 != nil {
				http.Error(w, "student_id OR booking_id should be a number", http.StatusBadRequest)
				return
			}
			var booked config.Booked
			booked.ID = uint(bkid)
			booked.StudentId = uint(id)
			result3 := config.Database.Where("id = ?", booked.ID).Find(&booked)
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

		} else {
			http.Error(w, "You are not authorised to cancel the booking", http.StatusBadRequest)
			fmt.Println("You are not authorised to cancel the booking")
		}

	}
}
