package controllers

import (
	"fmt"
	"net/http"

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

func Admin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		//e := "rak@toppr.com"

		fmt.Println(r.FormValue("uid"))

		//	uid =
		//fmt.Println("r.FormValue()")

		uid := r.FormValue("uid")

		idtodisable := r.FormValue("idtodisable")
		//request to disable a student or teacher

		u, ok := models.FindUserFromId(idtodisable)

		if !ok {
			http.Error(w, "no user found", http.StatusForbidden)
			fmt.Println("no user found")
			return
		}

		if u.Identity < "2" {

			if uid >= "2" {
				//if user iddentity is> 2  means that active user is an admin or super admin & has rights to make any user inactive
				u := models.MakeInactive(idtodisable)
				fmt.Print(u)

			} else {
				fmt.Print("You do not have the rights to make user inactive")
			}
		} else if u.Identity == "2" { //request to disable admin do only super admin can do so

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
			http.Error(w, "username does not exits", http.StatusForbidden)
			fmt.Println("Logined Failed")
			return
		}

		d, err := models.IsDisabled(u)
		if d == true {
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

// Logout method to call when the user signed out of the application.
func Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
		})
		fmt.Println("LogOut Successfully")
		w.WriteHeader(http.StatusOK)

		// http.Redirect(w, r, "/", http.StatusOK)
	}

}
