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

	if !valid(email) {
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

func MakeInactive(id string) models.User {

	user := models.User{}

	user, ok := FindUserFromId(id)

	identity := user.Identity

	if identity == 1 { ///he is a student

		var booked []models.Booked

		result1 := Database.Where("student_id= ?", user.ID).Find(&booked)
		if result1.Error != nil {

			return user
		}

		for i := range booked {

			result2 := Database.Model(&models.Slot{}).Where("id = ? ", booked[i].SlotId).Update("is_booked", 0)
			if result2.Error != nil {
				return user
			}

			result3 := Database.Delete(&booked[i])

			if result3.Error != nil {
				return user
			}

		}

	} else if identity == 0 { ///he is teacher

		var slots []models.Slot

		result1 := Database.Where("teacher_id= ?", user.ID).Find(&slots)
		if result1.Error != nil {

			return user
		}

		for i := range slots {
			result2 := Database.Model(&models.Booked{}).Where("slot_id = ? ", slots[i].ID).Delete(&models.Booked{})
			if result2.Error != nil {

				return user
			}

			result3 := Database.Delete(&slots[i])

			if result3.Error != nil {
				return user

			}
		}

		if !ok {
			fmt.Println("Logined Failed")
			return user
		}
	} else {
		fmt.Println("he is an admin/superadmin")
	}

	Database.Model(&user).Update("isdisabled", 1)

	return user
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
	return err == nil
}
