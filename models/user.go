package models

import (
	"errors"
	"net/http"

	"github.com/dipesh-toppr/bfsbeapp/config"
	"golang.org/x/crypto/bcrypt"
)

// User object handles information about application's registered users.
type User struct {
	ID         uint64
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

func IsDisabled(u User) (bool, error) {

	if u.Isdisabled == "1" {
		return true, errors.New("User is disabled  by admin")
	}

	return false, nil
}

// ValidatePassword validates the input password against the one in the database.
func (u *User) ValidatePassword(p string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))

	return err == nil

}
