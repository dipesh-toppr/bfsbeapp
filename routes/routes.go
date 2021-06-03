package routes

import (
	"log"
	"net/http"

	"github.com/dipesh-toppr/bfsbeapp/controllers"
	"github.com/dipesh-toppr/bfsbeapp/token"
	"github.com/gorilla/mux"
)

// LoadRoutes handles routes to pages of the application.
func LoadRoutes() {
	r := mux.NewRouter()
	// Index or main page.
	r.HandleFunc("/", index)

	// User related route(s)
	r.HandleFunc("/signup", controllers.Signup)
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)

	// welcome page
	r.HandleFunc("/welcome", welcome)
	r.HandleFunc("/search-teacher", controllers.SearchTeacher)
	log.Fatal(http.ListenAndServe(":8080", r))

}

// just check index page
func index(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("go ahead, its ok !"))
	w.WriteHeader(http.StatusOK)
	// http.Redirect(w, r, "/", http.StatusOK)

}

// try welcome api for fun !
func welcome(w http.ResponseWriter, r *http.Request) {

	e := token.Parsetoken(w, r)
	if e != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	// http.Redirect(w, r, "/", http.StatusOK)
}
