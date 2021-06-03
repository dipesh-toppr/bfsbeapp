package controllers

import (
	"fmt"
	"net/http"

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

// search teahcer for specific timing and subject

func SearchTeacher(w http.ResponseWriter, r *http.Request) {
	sub := r.URL.Query()["sub"][0]
	tim := r.URL.Query()["time"][0]
	fmt.Println(sub + " " + tim)
	var users []models.User
	////find the users
	// config.Database.Find(&users)
	// for _, val := range users {
	// 	fmt.Println()
	// }
	//update user details
	// config.Database.Model(&models.User{}).Where("id = ? ", 6).Update("identity", "1")

	///find user  with conditions
	// config.Database.Find(&users, "id = ?", 6)
	// for _, val := range users {
	// 	fmt.Println(val)
	// }
	///delete user with
	// var user models.User
	// config.Database.Where("id = ?", 6).Delete(&user)

	result := config.Database.Find(&users, "id = ?", 6)
	if len(users) == 0 {
		fmt.Println("not found")
	}
	fmt.Println(result.Error)
}

// Logout method to call when the user signed out of the application.
func Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		fmt.Println("LogOut Successfully")
		w.WriteHeader(http.StatusOK)

		// http.Redirect(w, r, "/", http.StatusOK)
	}

}
