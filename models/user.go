package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int
	Email      string
	Password   string
	Firstname  string
	Lastname   string
	Identity   string
	Isdisabled string
}

// SaveUser create new user entry.
func SaveUser(r *http.Request) (User, error) {

	// Get form values and validate
	u, err := validateUserForm(r)

	if err != nil {
		return u, err
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)

	if err != nil {
		return u, errors.New("the provided password is not valid")
	}

	u.Password = string(bs)

	if config.Database.Create(&u).Error != nil {
		return u, errors.New("unable to process registration")
	}
	return u, nil
}

// ValidateForm validates the submitted form for registration
func validateUserForm(r *http.Request) (User, error) {

	u := User{}
	e := r.FormValue("email")
	p := r.FormValue("password")
	cp := r.FormValue("cpassword")
	f := r.FormValue("firstname")
	l := r.FormValue("lastname")

	i := r.FormValue("identity")
	fmt.Println("hii")
	fmt.Println(f)

	if p != cp {
		return u, errors.New("password does not match")
	}

	if e == "" || p == "" || cp == "" || i == "" {
		return u, errors.New("fields cannot be empty")
	}

	_, err := CheckUser(e)

	if err != nil {
		return u, err
	}

	u.Email = e
	u.Firstname = f
	u.Lastname = l
	u.Password = p
	u.Identity = i
	u.Isdisabled = "0"

	return u, nil

}

// CheckUser looks for existing user using email
func CheckUser(email string) (User, error) {

	usr, ok := FindUser(email)

	if ok {
		return usr, errors.New("email is already taken")
	}

	return usr, nil

}

// FindUser looks for registerd user by email.
func FindUser(email string) (User, bool) {

	u := User{}

	if config.Database.Where(&User{Email: email}).Find(&u).Error != nil {
		return u, false
	}

	return u, true

}

func MakeInactive(id string) User {

	u := User{}

	u, ok := FindUserFromId(id)

	iden := u.Identity
	fmt.Println(iden)
	fmt.Println(u)

	if iden == "1" { ///he is a student

		var booked []config.Booked

		result1 := config.Database.Where("student_id= ?", u.ID).Find(&booked)
		if result1.Error != nil {

			return u
		}

		fmt.Println(booked)

		//result2 := config.Database.Where("slot_id= ?",booked.SlotId).Find(&booked)
		for i := range booked {

			result2 := config.Database.Model(&config.Slot{}).Where("id = ? ", booked[i].SlotId).Update("is_booked", 0)
			if result2.Error != nil {
				//	http.Error(w, result2.Error.Error(), http.StatusInternalServerError)
				return u
			}

			result3 := config.Database.Delete(&booked[i])

			if result3.Error != nil {
				//http.Error(w, result2.Error.Error(), http.StatusInternalServerError)
				return u
			}

		}

	} else if iden == "0" { ///he is teacher

		var slots []config.Slot

		result1 := config.Database.Where("teacher_id= ?", u.ID).Find(&slots)
		if result1.Error != nil {

			return u
		}

		fmt.Println(slots)

		//result2 := config.Database.Where("slot_id= ?",booked.SlotId).Find(&booked)
		for i := range slots {
			result2 := config.Database.Model(&config.Booked{}).Where("slot_id = ? ", slots[i].ID).Delete(&config.Booked{})
			if result2.Error != nil {

				return u
			}

			result3 := config.Database.Delete(&slots[i])

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

	//db.Model(&u).Update("isdisabled", "0")

	config.Database.Model(&u).Update("isdisabled", "1")

	return u
}

func FindUserFromId(id string) (User, bool) {

	u := User{}

	i, _ := strconv.Atoi(id)
	if config.Database.Where(&User{ID: i}).Find(&u).Error != nil {
		return u, false
	}

	return u, true

}

func IsDisabled(u User) (bool, error) {

	fmt.Print(u.Isdisabled)

	if u.Isdisabled == "1" {
		return true, errors.New("user is disabled  by admin")
	}

	return false, nil
}

// ValidatePassword validates the input password against the one in the database.
func (u *User) ValidatePassword(p string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))

	return err == nil

}

//find the type of user
func UserType(uid uint) string {
	var user User
	config.Database.Where("id = ?", uid).Find(&user)
	return user.Identity
}
