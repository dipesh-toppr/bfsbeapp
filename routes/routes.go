package routes

import (
	"log"
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

	//slot related routes
	http.HandleFunc("/addSlot", controllers.AddSlot)
	http.HandleFunc("/getUserSlots", controllers.GetUserSlots)
	http.HandleFunc("/updateSlot", controllers.UpdateSlot)
	http.HandleFunc("/deleteSlot", controllers.DeleteSlot)
	http.HandleFunc("/getUniqueSlots", controllers.GetUniqueSlots)

	//booking related routes
	http.HandleFunc("/search-teacher", controllers.SearchTeacher)
	http.HandleFunc("/delete-booking", controllers.DeleteBooking)
	http.HandleFunc("/read-booking", controllers.ReadBooking)

	//admin disable
	http.HandleFunc("/admin", controllers.Admin)

	http.HandleFunc("/admindeletebooking", controllers.AdminDeleteBooking)

	http.HandleFunc("/readallbookings", controllers.ReadAllBookings)

	http.HandleFunc("/readallstudents", controllers.ReadAllStudents)

	http.HandleFunc("/readallteachers", controllers.ReadAllTeachers)

	// welcome page
	// http.HandleFunc("/welcome", welcome)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// just check index page
func index(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("go ahead, its ok !"))
	w.WriteHeader(http.StatusOK)
	// http.Redirect(w, r, "/", http.StatusOK)

}

// try welcome api for fun !
// func welcome(w http.ResponseWriter, r *http.Request) {

// 	e, mail := token.Parsetoken(w, r)
// 	fmt.Printf(mail)
// 	if e != nil {
// 		http.Redirect(w, r, "/", http.StatusUnauthorized)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	// http.Redirect(w, r, "/", http.StatusOK)
// }
