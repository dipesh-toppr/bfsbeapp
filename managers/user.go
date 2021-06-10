package managers

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/dipesh-toppr/bfsbeapp/models"
	"golang.org/x/crypto/bcrypt"
)

// check validity of user form form signup api
func ValidateUserFormSignup(request *http.Request, response http.ResponseWriter) (map[string]interface{}, error) {

	validparams := map[string]interface{}{}

	email := trim(request.FormValue("email"))
	password := trim(request.FormValue("password"))
	confirmpassword := trim(request.FormValue("cpassword"))
	firstname := trim(request.FormValue("firstname"))
	lastname := trim(request.FormValue("lastname"))
	identity := trim(request.FormValue("identity"))

	if email == "" || password == "" || firstname == "" {
		return validparams, errors.New("fields cannot be empty")
	}

	if password != confirmpassword {
		return validparams, errors.New("password does not match")
	}

	if identity != models.IDENTITY["teacher"] && identity != models.IDENTITY["student"] {
		return validparams, errors.New("wrong identity passed")
	}

	if !valid(email) {
		return validparams, errors.New("email is not valid")
	}

	_, exists := FindUser(email)

	if exists {
		return validparams, errors.New("email is already taken")
	}

	validparams["Email"] = email
	validparams["Firstname"] = firstname
	validparams["Lastname"] = lastname
	validparams["Password"] = password
	validparams["Identity"], _ = strconv.Atoi(identity)
	validparams["Isdisabled"] = false

	return validparams, nil
}

// looks for registerd user by email.
func FindUser(email string) (models.User, bool) {

	user := models.User{}

	if Database.Where(&models.User{Email: email}).Find(&user).Error != nil {
		return user, false
	}

	return user, true

}

// SaveUser create new user entry.
func SaveUser(validparams map[string]interface{}) (models.User, error) {

	bctyptpassword, err := bcrypt.GenerateFromPassword([]byte(validparams["Password"].(string)), bcrypt.MinCost)

	if err != nil {
		return models.User{}, errors.New("the provided password is not valid")
	}

	validparams["Password"] = string(bctyptpassword)
	user, err := CreateUser(validparams)

	if err != nil {
		return user, err
	}
	return user, nil
}

// create new user
func CreateUser(validparams map[string]interface{}) (models.User, error) {

	user := models.User{}
	user.Email = validparams["Email"].(string)
	user.Firstname = validparams["Firstname"].(string)
	user.Lastname = validparams["Lastname"].(string)
	user.Password = validparams["Password"].(string)
	user.Identity = validparams["Identity"].(int)
	user.Isdisabled = validparams["Isdisabled"].(bool)

	if Database.Create(&user).Error != nil {
		return user, errors.New("unable to process registration")
	}

	return user, nil
}

// check validity of user form form login api
func ValidateUserFormLogin(request *http.Request, response http.ResponseWriter) (map[string]interface{}, error) {

	validparams := map[string]interface{}{}

	email := trim(request.FormValue("email"))
	password := trim(request.FormValue("password"))

	if email == "" || password == "" {
		return validparams, errors.New("fields cannot be empty")
	}

	if !valid(email) {
		return validparams, errors.New("email is not valid")
	}

	validparams["Email"] = email
	validparams["Password"] = password

	return validparams, nil
}

// ValidatePassword validates the input password against the one in the database.
func ValidatePassword(user models.User, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	return err == nil

}

// trim if any space in input
func trim(params string) string {
	return strings.TrimSpace(params)
}

// find user by id
func FindUserFromId(id string) (models.User, bool) {

	user := models.User{}

	i, _ := strconv.Atoi(id)
	if Database.Where(&models.User{ID: i}).Find(&user).Error != nil {
		return user, false
	}

	return user, true

}

// check if user is disabled by admin or not
func IsDisabled(user models.User) (bool, error) {

	if user.Isdisabled {
		return true, errors.New("user is disabled  by admin")
	}

	return false, nil
}

//find the type of user
func UserType(uid uint) int {

	var user models.User
	Database.Where("id = ?", uid).Find(&user)

	return user.Identity
}

// check validity of the email
func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
