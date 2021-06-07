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

// Signup allows the user to create an account.
func Signup(w http.ResponseWriter, r *http.Request) {

	// var u models.User
	// process form submission
	if r.Method == http.MethodPost {
		var u models.User
		u, err := models.SaveUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("SignUp Failed")
			return
		}

		fmt.Println(u)
		fmt.Println("SignUp Successfully")

		// add token to cookies
		token.Createtoken(u, w)
		w.WriteHeader(http.StatusOK)

		// redirect
		// http.Redirect(w, r, "/", http.StatusOK)
		return
	}
}

// Login allows registered user to access the application.
func Login(w http.ResponseWriter, r *http.Request) {

	// var u models.User
	// process form submission
	if r.Method == http.MethodPost {

		p := r.FormValue("password")
		e := r.FormValue("email")

		// check if the user exists
		u, ok := models.FindUser(e)
		if !ok {
			http.Error(w, "username and/or password do not match", http.StatusForbidden)
			fmt.Println("Logined Failed")
			return
		}

		d, err := models.IsDisabled(u)
		if d {
			http.Error(w, err.Error(), http.StatusForbidden)
			fmt.Println("user is disabled by admin....")
			return
		}

		if !u.ValidatePassword(p) {
			http.Error(w, "username and/or password do not match", http.StatusForbidden)
			fmt.Println("Logined Failed")
			return
		}

		fmt.Println("Logined Successfully")

		// add token to cookies
		token.Createtoken(u, w)
		w.WriteHeader(http.StatusOK)

		// http.Redirect(w, r, "/", http.StatusOK)
		return
	}

}

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

			} else {
				fmt.Print("You do not have the rights to make user inactive")
			}
		} else if utodisable.Identity == "2" { //request to disable admin do only super admin can do so

			if uid == "3" { //identity  of superadmin  kept 3
				u := models.MakeInactive(idtodisable)
				fmt.Print(u)
			} else {
				fmt.Print("You do not have the rights to make admin inactive")
			}
		}
		// http.Redirect(w, r, "/", http.StatusOK)
	}

}

// search teahcer for specific timing

func SearchTeacher(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)
		if e != nil {
			http.Error(w, "unauthorized request", http.StatusBadRequest)
			return
		}
		tim, err := models.ValidateTime(r)
		print(tim, " ", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

		} else {
			fmt.Println("You are not authorised to cancel the booking")
		}

	}
}

// Logout method to call when the user signed out of the application.
func Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	fmt.Println("LogOut Successfully")
	w.WriteHeader(http.StatusOK)

}
