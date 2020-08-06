package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type authenticationData struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login enpoint hit")
	var creds credentials
	var loginResponseData authenticationData

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		//Send back 400
		loginResponseData = authenticationData{
			Username: "Bad Request",
		}
		sendFailResponse(w, loginResponseData)
		return
	}

	usr, err := obtainUserFromSliceByUsername(creds.Username)
	if err != nil || usr.Password != creds.Password {
		loginResponseData = authenticationData{
			Username: "Invalid login",
		}
		unauthorizedResponse(w, loginResponseData)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claim := &claims{
		ID: usr.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		//internal server error response 500 this needs to be changed
		sendFailResponse(w, "Internal failure")
		return
	}

	//Send username and id to reduce api calls
	loginResponseData = authenticationData{
		Username: creds.Username,
		ID:       usr.ID,
		Token:    tokenString,
	}

	sendSuccessResponse(w, loginResponseData)
}

func validateToken(w http.ResponseWriter, r *http.Request) bool {
	var authResponse authenticationData
	authTokenString := r.Header.Get("Authorization")

	splitToken := strings.Split(authTokenString, "Bearer")
	if len(splitToken) != 2 {
		authResponse = authenticationData{
			Token: "Invalid Request",
		}
		//return 400
		sendFailResponse(w, authResponse)
		return false
	}
	jwtToken := strings.TrimSpace(splitToken[1])

	//obtain token from authorization header
	claim := &claims{}

	tkn, err := jwt.ParseWithClaims(jwtToken, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			authResponse = authenticationData{
				Token: "Unauthroized Request",
			}
			//return 403
			unauthorizedResponse(w, authResponse)
			return false
		}
		authResponse = authenticationData{
			Token: "Invalid Request",
		}
		//return 400
		sendFailResponse(w, authResponse)
		return false
	}
	if !tkn.Valid {
		authResponse = authenticationData{
			Token: "Unauthroized Request",
		}
		//return 400
		unauthorizedResponse(w, authResponse)
		return false
	}
	return true
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	var authResponse authenticationData
	authTokenString := r.Header.Get("Authorization")

	splitToken := strings.Split(authTokenString, "Bearer")
	if len(splitToken) != 2 {
		authResponse = authenticationData{
			Token: "Invalid Request",
		}
		//return 400
		sendFailResponse(w, authResponse)
		return
	}
	jwtToken := strings.TrimSpace(splitToken[1])

	//obtain token from authorization header
	claim := &claims{}

	tkn, err := jwt.ParseWithClaims(jwtToken, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			authResponse = authenticationData{
				Token: "Unauthroized Request",
			}
			//return 403
			unauthorizedResponse(w, authResponse)
			return
		}
		authResponse = authenticationData{
			Token: "Invalid Request",
		}
		//return 400
		sendFailResponse(w, authResponse)
		return
	}
	if !tkn.Valid {
		authResponse = authenticationData{
			Token: "Unauthroized Request",
		}
		//return 400
		unauthorizedResponse(w, authResponse)
		return
	}

	//checks if its earlier than 30 seconds
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		authResponse = authenticationData{
			Token: "Invalid Request",
		}
		//return 400
		sendFailResponse(w, authResponse)
		return
	}

	//create new token
	expirationTime := time.Now().Add(5 * time.Minute)
	claim.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		//send 500...
		sendFailResponse(w, "Somthing went wrong....")
		return
	}

	authResponse = authenticationData{
		Token: tokenString,
	}
	sendSuccessResponse(w, authResponse)
}

func bootstrapFromToken(w http.ResponseWriter, r *http.Request) {
	var authResponse authenticationData
	authTokenString := r.Header.Get("Authorization")

	splitToken := strings.Split(authTokenString, "Bearer")
	if len(splitToken) != 2 {
		authResponse = authenticationData{
			Token: "Invalid Request",
		}
		//return 400
		sendFailResponse(w, authResponse)
		return
	}
	jwtToken := strings.TrimSpace(splitToken[1])

	//obtain token from authorization header
	claim := &claims{}

	tkn, err := jwt.ParseWithClaims(jwtToken, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			authResponse = authenticationData{
				Token: "Unauthorized Request",
			}
			//return 403
			unauthorizedResponse(w, authResponse)
			return
		}
		authResponse = authenticationData{
			Token: "Invalid Request",
		}
		//return 400
		sendFailResponse(w, authResponse)
		return
	}
	if !tkn.Valid {
		authResponse = authenticationData{
			Token: "Unauthroized Request",
		}
		//return 400
		unauthorizedResponse(w, authResponse)
		return
	}

	//Should we attempt to refresh to token if its too old??

	//When we bootstrap we want to send the username and id to reduce api calls
	usr, _ := obtainUserFromSlice(claim.ID)
	authResponse = authenticationData{
		ID:       claim.ID,
		Username: usr.Username,
	}
	sendSuccessResponse(w, authResponse)
}
