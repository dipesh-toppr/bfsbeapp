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

		// add token to cookies
		token.Createtoken(u, w)
		w.WriteHeader(http.StatusOK)

		fmt.Println(u)
		w.Write([]byte("SignUp Successful"))
		fmt.Println("SignUp Successful")

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

		d, _ := models.IsDisabled(u)
		if d {
			// http.Error(w, err.Error(), http.StatusForbidden)
			http.Error(w, "user is disabled by admin....", http.StatusForbidden)
			fmt.Println("user is disabled by admin....")
			return
		}

		if !u.ValidatePassword(p) {
			http.Error(w, "username and/or password do not match", http.StatusForbidden)
			fmt.Println("Logined Failed")
			return
		}

		// add token to cookies
		token.Createtoken(u, w)
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Login Successful"))
		fmt.Println("Login Successful")

		// http.Redirect(w, r, "/", http.StatusOK)
		return
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
	w.Write([]byte("LogOut Successful"))
	fmt.Println("LogOut Successful")
	w.WriteHeader(http.StatusOK)

}
