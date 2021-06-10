package managers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/models"
	"github.com/dipesh-toppr/bfsbeapp/token"
)

func AdminDeleteBook(w http.ResponseWriter, r *http.Request) {

	id, e := token.Parsetoken(w, r) //finding active user mail and if he is logged in

	// print(e, "  ", mail) //for debugging

	if e != nil {
		//	http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	user, ok := FindUserFromId(strconv.Itoa(int(id))) //find the user from id

	if !ok {
		http.Error(w, "no user found", http.StatusNotFound)
		fmt.Println("no user found")
		return
	}

	uid := user.Identity

	if uid > 1 {

		bid := r.FormValue("bid")
		bkid, err2 := strconv.Atoi(bid)
		if err2 != nil {
			http.Error(w, "student_id OR booking_id should be a number", http.StatusBadRequest)
			return
		}
		var booked models.Booked
		booked.ID = uint(bkid)
		booked.StudentId = uint(id)
		result3 := Database.Where("id = ?", booked.ID).Find(&booked)
		slot := booked.SlotId
		if result3.Error != nil {
			http.Error(w, "Invalid booking ID", http.StatusBadRequest)
			return
		}
		result1 := Database.Where("id = ?", booked.ID).Delete(&booked)
		if result1.Error != nil {
			http.Error(w, result1.Error.Error(), http.StatusInternalServerError)
			return
		}
		result2 := Database.Model(&models.Slot{}).Where("id = ? ", slot).Update("is_booked", 0)
		if result2.Error != nil {
			http.Error(w, result2.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Booking Deleted!"))
		w.WriteHeader(http.StatusOK)
		return

	} else {
		http.Error(w, "You are not authorised to cancel the booking", http.StatusBadRequest)
		fmt.Println("You are not authorised to cancel the booking")
	}

}

func MakeInactive(id string) models.User {

	u := models.User{}

	u, ok := FindUserFromId(id)

	iden := u.Identity
	fmt.Println(iden)
	fmt.Println(u)

	if iden == 1 { ///he is a student

		var booked []models.Booked

		result1 := Database.Where("student_id= ?", u.ID).Find(&booked)
		if result1.Error != nil {

			return u
		}

		fmt.Println("booked", booked)
		slotids := []int{}

		for i := range booked {
			slotids = append(slotids, int(booked[i].SlotId))
		}
		fmt.Println("slotid", slotids)
		//result2 := Database.Where("slot_id= ?",booked.SlotId).Find(&booked)
		//
		for i := range slotids {

			result2 := Database.Model(&models.Slot{}).Where("id = ?", slotids[i]).Update("is_booked", false)

			//Database.Raw("UPDATE slots SET is_booked = ? WHERE id IN ? ", false, &slotids)

			if result2.Error != nil {
				//	http.Error(w, result2.Error.Error(), http.StatusInternalServerError)
				return u
			}

			result3 := Database.Delete(&booked[i])

			if result3.Error != nil {
				//http.Error(w, result3.Error(), http.StatusInternalServerError)
				return u
			}
		}

	} else if iden == 0 { ///he is teacher

		var slots []models.Slot

		result1 := Database.Where("teacher_id= ?", u.ID).Find(&slots)
		if result1.Error != nil {

			return u
		}

		fmt.Println(slots)

		ids := []int{}

		for i := range slots {
			ids = append(ids, int(slots[i].ID))
		}

		fmt.Println("ids", ids)

		//var bookingstodel []models.Booked

		//result2 := Database.Where("slot_id= ?",booked.SlotId).Find(&booked)
		for i := range slots {

			result2 := Database.Model(&models.Booked{}).Where("slot_id = ? ", slots[i].ID).Delete(&models.Booked{})
			//	result2 := Database.Model(&models.Booked{}).Where("slot_id IN ?", ids).Find(&bookingstodel)
			if result2.Error != nil {
				return u
			}

			result3 := Database.Delete(&slots[i])
			//result4 := Database.Delete(&bookingstodel)
			if result3.Error != nil {
				//http.Error(w, result2.Error.Error(), http.StatusInternalServerError)
				return u

			}

		}

		if !ok {
			//http.Error(w, "username does not exits", http.StatusForbidden)
			fmt.Println("Logined Failed")
			return u
		}
	} else {

		fmt.Println("he is an admin/superadmin")
	}

	//db.Model(&u).Update("isdisabled", 0)

	Database.Model(&u).Update("isdisabled", true)

	return u
}

func ReadStudents(r *http.Request) ([]models.User, bool) {

	var stud []models.User
	result := Database.Where("identity = ?", uint(1)).Find(&stud)
	if result.Error != nil {
		return stud, false
	}
	fmt.Print(stud)
	return stud, true
}

func ReadTeachers(r *http.Request) ([]models.User, bool) {

	var teach []models.User
	result := Database.Where("identity = ?", uint(0)).Find(&teach)
	if result.Error != nil {
		return teach, false
	}
	fmt.Print(teach)
	return teach, true
}

func ReadAdminBooked(r *http.Request) ([]models.Booked, bool) {

	var booked []models.Booked
	result := Database.Find(&booked)
	if result.Error != nil {
		return booked, false
	}

	return booked, true
}
