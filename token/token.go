package token

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dipesh-toppr/bfsbeapp/models"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Email  string `json:"email"`
	UserId int    `json:"user_id"`
	jwt.StandardClaims
}

func init() {}

func Createtoken(u models.User, w http.ResponseWriter) error {
	// Create the JWT claims, which includes the EmailId.
	claims := &Claims{
		Email:          u.Email,
		UserId:         (u.ID),
		StandardClaims: jwt.StandardClaims{},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
	})

	return nil
}

func Parsetoken(w http.ResponseWriter, r *http.Request) (int, error) {

	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return 0, err
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return 0, err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return 0, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return 0, err
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("TOKEN NOT VALID")
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	// w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Email)))

	return claims.UserId, nil
}
