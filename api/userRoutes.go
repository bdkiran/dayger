package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	// err := validateToken(r)
	// if err != nil {
	// 	unauthorizedResponse(w, "Unauthroized")
	// 	return
	// }

	variables := mux.Vars(r)
	userID := variables["id"]

	userToReturn, err := obtainUserFromSlice(userID)
	if err != nil {
		sendFailResponse(w, userToReturn)
		return
	}
	sendSuccessResponse(w, userToReturn)
}

//Create user
//update user
//Delete user
