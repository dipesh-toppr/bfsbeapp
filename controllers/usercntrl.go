package controllers

import (
	"net/http"

	"github.com/dipesh-toppr/bfsbeapp/managers"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

// Signup allows the user to create an account.
func Signup(response http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodPost {

		validparams, err := managers.ValidateUserFormSignup(request, response)

		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := managers.SaveUser(validparams)

		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}

		err = token.Createtoken(user, response)

		if err != nil {
			return
		}

		response.WriteHeader(http.StatusOK)
		response.Write([]byte("SignUp Successful"))
		return
	}
}

// Login allows registered user to access the application.
func Login(response http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodPost {

		validparams, err := managers.ValidateUserFormLogin(request, response)

		if err != nil {
			http.Error(response, err.Error(), http.StatusNotFound)
			return
		}

		// check if the user exists
		user, ok := managers.FindUser(validparams["Email"].(string))
		if !ok {
			http.Error(response, "username and/or password do not match", http.StatusNotFound)
			return
		}

		disable, err := managers.IsDisabled(user)
		if disable {
			http.Error(response, err.Error(), http.StatusForbidden)
			return
		}

		if !managers.ValidatePassword(user, validparams["Password"].(string)) {
			http.Error(response, "username and/or password do not match", http.StatusForbidden)
			return
		}

		err = token.Createtoken(user, response)

		if err != nil {
			return
		}

		response.WriteHeader(http.StatusOK)
		response.Write([]byte("Login Successful"))
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
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("LogOut Successful"))

}
