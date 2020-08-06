package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Returns an experience object using the experience id
func getExperienceByID(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	userID := variables["id"]
	experience, err := obtainExperienceFromSliceByID(userID)
	if err != nil {
		sendFailResponse(w, experience)
		return
	}

	sendSuccessResponse(w, experience)
}

//Returns a list of all the experiences using the user id
func getExperiencesByUserID(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	userID := variables["id"]

	experiences, err := obtainExperiencesFromSliceByUserID(userID)
	if err != nil {
		sendFailResponse(w, experiences)
		return
	}
	sendSuccessResponse(w, experiences)
}

//Creates a new experience
func createNewExperience(w http.ResponseWriter, r *http.Request) {
	var ex experience
	decoder := json.NewDecoder(r.Body)
	//currently sending the id, but in the future needs to be generated server/db side to avoid collisions
	err := decoder.Decode(&ex)
	if err != nil {
		sendFailResponse(w, "Invalid request. Not JSON encoded.")
		return
	}

	experience, err := addExperienceToSlice(ex)
	if err != nil {
		sendFailResponse(w, experience)
		return
	}

	sendSuccessResponse(w, experience)
}

//Removes an experience
func deleteExperience(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	userID := variables["id"]

	experience, err := removeExperienceFromSliceByID(userID)
	if err != nil {
		sendFailResponse(w, experience)
		return
	}

	sendSuccessResponse(w, experience)
}

//Updates an experience
func updateExperience(w http.ResponseWriter, r *http.Request) {
	var ex experience
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ex)
	if err != nil {
		sendFailResponse(w, "Invalid request. Not JSON encoded.")
		return
	}
	experience, err := updateExperienceInSlice(ex)
	if err != nil {
		sendFailResponse(w, experience)
		return
	}

	sendSuccessResponse(w, experience)
}
