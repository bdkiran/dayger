package api

import (
	"errors"
	"fmt"
	"log"
)

//Do I want to keep the omitempty portion, this is used for formating error responses
type user struct {
	Username string `json:"username,omitempty"`
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

/*Functions interacting with our user static test data*/
var testUsers = []user{
	user{"bcastle", "0", "password1"},
	user{"samBerry", "1", "password1"},
}

//Obtains a user from the hardcoded slice
func obtainUserFromSlice(id string) (user, error) {
	log.Printf("Looking for user: %s", id)
	for _, us := range testUsers {
		if us.ID == id {
			return us, nil
		}
	}
	errorMessage := fmt.Sprintf("Unable to find a user with the id: %s", id)
	matchingUser := user{"", errorMessage, ""}
	return matchingUser, errors.New(errorMessage)
}

//Obtains a user from the hardcoded slice
func obtainUserFromSliceByUsername(username string) (user, error) {
	log.Printf("Looking for user: %s", username)
	for _, us := range testUsers {
		if us.Username == username {
			return us, nil
		}
	}
	errorMessage := fmt.Sprintf("Unable to find a user with the username: %s", username)
	matchingUser := user{"", errorMessage, ""}
	return matchingUser, errors.New(errorMessage)
}

/*Functions interacting with our experience static test data*/

//Do i want to keep the omitempty portion, this is used for formating error responses
type experience struct {
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Time        string `json:"time,omitempty"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
	UserID      string `json:"userID,omitempty"`
}

var testExperiences = []experience{
	experience{"Awsome Throwdown", "0", "07-07-2020", "Holywood", "Awsome party bring chicks", "1"},
	experience{"Brunch and MOMOS", "1", "07-07-2020", "NYC", "Costume themed, dress up like wild, wild west", "0"},
}

func obtainExperienceFromSliceByID(id string) (experience, error) {
	log.Printf("Looking for Experience: %s", id)
	for _, ex := range testExperiences {
		if ex.ID == id {
			return ex, nil
		}
	}
	errorMessage := fmt.Sprintf("Unable to find an experience with the id: %s", id)
	returnExperience := experience{"", errorMessage, "", "", "", ""}
	return returnExperience, errors.New(errorMessage)
}

func obtainExperiencesFromSliceByUserID(userID string) ([]experience, error) {
	log.Printf("Looking for experience under the userID: %s", userID)
	var userExperiences []experience
	for _, ex := range testExperiences {
		if ex.UserID == userID {
			userExperiences = append(userExperiences, ex)
		}
	}
	//Will just return an empty array if there is no matching values...
	// Unsure of what errors are needed to bubble up.
	// errorMessage := fmt.Sprintf("Unable to find an experience with the user id: %s", userID)
	// returnExperience := experience{"", "", "", "", "", errorMessage}
	// return returnExperience, errors.New(errorMessage)
	return userExperiences, nil
}

func addExperienceToSlice(ex experience) (experience, error) {
	testExperiences = append(testExperiences, ex)
	return ex, nil
}

func removeExperienceFromSliceByID(id string) (experience, error) {
	var exp experience
	log.Printf("Attempting to delete experience: %s", id)
	for i, ev := range testExperiences {
		if ev.ID == id {
			//Mutate the gloabl array
			testExperiences[i] = testExperiences[len(testExperiences)-1]
			testExperiences = testExperiences[:len(testExperiences)-1]
			//returnString := "Succussfully deleted experience: " + id
			return exp, nil
		}
	}
	//Not sure if this message should fall under this catagory
	errorMessage := fmt.Sprintf("Unable to find an experience with the id: %s", id)
	exp = experience{"", errorMessage, "", "", "", ""}
	return exp, errors.New(errorMessage)
}

func updateExperienceInSlice(ex experience) (experience, error) {
	log.Printf("Attempting to update experience: %s", ex.ID)
	for i, ev := range testExperiences {
		if ev.ID == ex.ID {
			testExperiences[i] = ex
			return ex, nil
		}
	}
	//Not sure if this message should fall under this catagory
	errorMessage := fmt.Sprintf("Unable to find an experience with the id: %s", ex.ID)
	returnExperience := experience{"", errorMessage, "", "", "", ""}
	return returnExperience, errors.New(errorMessage)
}
