package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/managers"
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

		user, ok := managers.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusBadRequest)
			return
		}

		uid := user.Identity //uid is identity of user ie stud,tech,admin,superadmin

		//fmt.Println("This is uid ", uid)

		//finding the user detials to check his/her role

		utodisable, ok := managers.FindUserFromId(idtodisable)

		if !ok {
			http.Error(w, "no user found", http.StatusNotFound)
			///	fmt.Println("no user found")
			return
		}

		if utodisable.Identity < 2 { ///means he is stud or teacher so can be made inactive my both admin and super admin

			if uid >= 2 {
				//if user iddentity is>= 2  means that active user is an admin or super admin & has rights to make any user inactive
				u := managers.MakeInactive(idtodisable)
				fmt.Print(u)
				w.Write([]byte("User disabled\n"))
				json.NewEncoder(w).Encode(u)
				return

			} else {
				http.Error(w, "You do not have the rights to make admin inactive", http.StatusUnauthorized)
				//fmt.Print("You do not have the rights to make user inactive")
			}
		} else if utodisable.Identity == 2 { //request to disable admin do only super admin can do so

			if uid == 3 { //identity  of superadmin  kept 3
				u := managers.MakeInactive(idtodisable)
				fmt.Print(u)
				w.Write([]byte("User disabled\n"))
				json.NewEncoder(w).Encode(u)
				return
			} else {
				http.Error(w, "You do not have the rights to make admin inactive", http.StatusUnauthorized)
			}
		}
		if utodisable.Identity == 3 {
			http.Error(w, "you cannot disable super admin", http.StatusUnauthorized)
		}
		// http.Redirect(w, r, "/", http.StatusOK)
	}

}

//admin read all bookings

func ReadAllBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)

		user, ok := managers.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusNotFound)
			fmt.Println("no user found")
			return
		}

		//uid := user.Identity
		fmt.Println(user)
		if e != nil || user.Identity < 2 {
			http.Error(w, "unauthorized request", http.StatusForbidden)
			return
		}

		slot, ok := managers.ReadAdminBooked(r)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(slot)
		w.WriteHeader(http.StatusOK)

	}
}

func ReadAllTeachers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)

		user, ok := managers.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusNotFound)
			fmt.Println("no user found")
			return
		}

		//uid := user.Identity
		fmt.Println(user)

		if e != nil || user.Identity < 2 {
			http.Error(w, "unauthorized request", http.StatusForbidden)
			return
		}

		teachers, ok := managers.ReadTeachers(r)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(teachers)

		w.WriteHeader(http.StatusOK)

	}
}

func ReadAllStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, e := token.Parsetoken(w, r)
		fmt.Println(id)

		user, ok := managers.FindUserFromId(strconv.Itoa(int(id))) //find the user from id

		if !ok {
			http.Error(w, "no user found", http.StatusNotFound)
			fmt.Println("no user found")
			return
		}

		//uid := user.Identity
		fmt.Println(user)

		if e != nil || user.Identity < 2 {
			http.Error(w, "unauthorized request", http.StatusForbidden)
			return
		}

		students, ok := managers.ReadStudents(r)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(students)

		w.WriteHeader(http.StatusOK)
	}
}

//admin delete booking

func AdminDeleteBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		managers.AdminDeleteBook(w, r)

	}
}
