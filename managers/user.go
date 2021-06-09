package managers

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/dipesh-toppr/bfsbeapp/models"
	"golang.org/x/crypto/bcrypt"
)

// SaveUser create new user entry.
func SaveUser(r *http.Request) (models.User, error) {

	// Get form values and validate
	user, err := validateUserForm(r)

	if err != nil {
		return user, err
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

	if err != nil {
		return user, errors.New("the provided password is not valid")
	}

	user.Password = string(bs)

	if Database.Create(&user).Error != nil {
		return user, errors.New("unable to process registration")
	}
	return user, nil
}

// ValidateForm validates the submitted form for registration
func validateUserForm(request *http.Request) (models.User, error) {

	user := models.User{}
	email := request.FormValue("email")
	password := request.FormValue("password")
	confirmpassword := request.FormValue("cpassword")
	firstname := request.FormValue("firstname")
	lastname := request.FormValue("lastname")

	identity := request.FormValue("identity")

	if password != confirmpassword {
		return user, errors.New("password does not match")
	}

	if email == "" || password == "" || confirmpassword == "" || identity == "" {
		return user, errors.New("fields cannot be empty")
	}

	if valid(email) == false {
		return user, errors.New("email is not valid")
	}
	_, err := CheckUser(email)

	if err != nil {
		return user, err
	}

	user.Email = email
	user.Firstname = firstname
	user.Lastname = lastname
	user.Password = password
	user.Identity, _ = strconv.Atoi(identity)
	user.Isdisabled = false

	return user, nil

}

// CheckUser looks for existing user using email
func CheckUser(email string) (models.User, error) {

	user, ok := FindUser(email)

	if ok {
		return user, errors.New("email is already taken")
	}

	return user, nil

}

// FindUser looks for registerd user by email.
func FindUser(email string) (models.User, bool) {

	user := models.User{}

	if Database.Where(&models.User{Email: email}).Find(&user).Error != nil {
		return user, false
	}

	return user, true

}

func FindUserFromId(id string) (models.User, bool) {

	user := models.User{}

	i, _ := strconv.Atoi(id)
	if Database.Where(&models.User{ID: i}).Find(&user).Error != nil {
		return user, false
	}

	return user, true

}

func IsDisabled(user models.User) (bool, error) {

	fmt.Print(user.Isdisabled)

	if user.Isdisabled {
		return true, errors.New("user is disabled  by admin")
	}

	return false, nil
}

// ValidatePassword validates the input password against the one in the database.
func ValidatePassword(user models.User, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	return err == nil

}

//find the type of user
func UserType(uid uint) int {
	var user models.User
	Database.Where("id = ?", uid).Find(&user)
	return user.Identity
}

// check validity of the email.
func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	if err == nil {
		return true
	} else {
		return false
	}
}
