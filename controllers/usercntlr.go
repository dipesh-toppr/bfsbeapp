package controllers

import (
	"fmt"
	"net/http"

	"github.com/dipesh-toppr/bfsbeapp/managers"
	"github.com/dipesh-toppr/bfsbeapp/models"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

// Signup allows the user to create an account.
func Signup(response http.ResponseWriter, request *http.Request) {

	// var u models.User
	// process form submission
	if request.Method == http.MethodPost {
		var user models.User
		user, err := managers.SaveUser(request)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			fmt.Println("SignUp Failed")
			return
		}

		// add token to cookies
		token.Createtoken(user, response)
		response.WriteHeader(http.StatusOK)

		fmt.Println(user)
		response.Write([]byte("SignUp Successful"))
		fmt.Println("SignUp Successful")

		// redirect
		// http.Redirect(w, r, "/", http.StatusOK)
		return
	}
}

// Login allows registered user to access the application.
func Login(response http.ResponseWriter, request *http.Request) {

	// var u models.User
	// process form submission
	if request.Method == http.MethodPost {

		password := request.FormValue("password")
		email := request.FormValue("email")

		// check if the user exists
		user, ok := managers.FindUser(email)
		if !ok {
			http.Error(response, "username and/or password do not match", http.StatusForbidden)
			fmt.Println("Logined Failed")
			return
		}

		disable, _ := managers.IsDisabled(user)
		if disable {
			// http.Error(w, err.Error(), http.StatusForbidden)
			http.Error(response, "user is disabled by admin....", http.StatusForbidden)
			fmt.Println("user is disabled by admin....")
			return
		}

		if !managers.ValidatePassword(user, password) {
			http.Error(response, "username and/or password do not match", http.StatusForbidden)
			fmt.Println("Logined Failed")
			return
		}

		// add token to cookies
		token.Createtoken(user, response)
		response.WriteHeader(http.StatusOK)

		response.Write([]byte("Login Successful"))
		fmt.Println("Login Successful")

		// http.Redirect(w, r, "/", http.StatusOK)
		return
	}

}

// Logout method to call when the user signed out of the application.
func Logout(response http.ResponseWriter, request *http.Request) {

	http.SetCookie(response, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	response.Write([]byte("LogOut Successful"))
	fmt.Println("LogOut Successful")
	response.WriteHeader(http.StatusOK)

}
