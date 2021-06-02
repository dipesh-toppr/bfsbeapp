package routes

import (
	"net/http"

	"github.com/dipesh-toppr/bfsbeapp/controllers"
)

// LoadRoutes handles routes to pages of the application.
func LoadRoutes() {

	// Index or main page.
	http.HandleFunc("/", index)

	// User related route(s)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)

	http.ListenAndServe(":8080", nil)

}

// Redirect to list of index page.
func index(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/", http.StatusOK)

}
