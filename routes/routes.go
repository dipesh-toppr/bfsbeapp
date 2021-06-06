package routes

import (
	"net/http"

	"github.com/dipesh-toppr/bfsbeapp/controllers"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

// LoadRoutes handles routes to pages of the application.
func LoadRoutes() {

	// Index or main page.
	http.HandleFunc("/", index)

	// User related route(s)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/addSlot", controllers.AddSlot)
	http.HandleFunc("/getUserSlots", controllers.GetUserSlots)
	http.HandleFunc("/updateSlot", controllers.UpdateSlot)
	http.HandleFunc("/deleteSlot", controllers.DeleteSlot)
	// welcome page
	http.HandleFunc("/welcome", welcome)

	http.ListenAndServe(":8080", nil)

}

// just check index page
func index(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("go ahead, its ok !"))
	w.WriteHeader(http.StatusOK)
	// http.Redirect(w, r, "/", http.StatusOK)

}

// try welcome api for fun !
func welcome(w http.ResponseWriter, r *http.Request) {

	e, _ := token.Parsetoken(w, r)
	if e != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	// http.Redirect(w, r, "/", http.StatusOK)
}
