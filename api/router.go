package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//NewRouter creates and reaturns a new routher that will register all routes
func NewRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")

	router.HandleFunc("/experience/user/{id}", getExperiencesByUserID).Methods("GET")
	router.HandleFunc("/experience/{id}", getExperienceByID).Methods("GET")
	router.HandleFunc("/experience", createNewExperience).Methods("POST")
	router.HandleFunc("/experience/{id}", deleteExperience).Methods("DELETE")
	router.HandleFunc("/experience", updateExperience).Methods("PUT")

	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/refresh", refreshToken).Methods("GET")
	router.HandleFunc("/bootstrap", bootstrapFromToken).Methods("GET")

	//Generate inviteID
	//Verify inviteID

	//Code for audio streaming
	songsDir := "/songs/"
	//router.PathPrefix(songsDir).Handler(http.StripPrefix(songsDir, http.FileServer(http.Dir("."+songsDir))))
	router.PathPrefix(songsDir).Handler(songFileHandler(songsDir))

	picutreDir := "/images/"
	router.PathPrefix(picutreDir).Handler(songFileHandler(picutreDir))

	router.Use(authMiddleWare)

	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	returnCors := handlers.CORS(headersOk, originsOk, methodsOk)(router)
	return returnCors
}

func songFileHandler(directory string) http.Handler {
	return http.StripPrefix(directory, http.FileServer(http.Dir("."+directory)))
}

func home(w http.ResponseWriter, r *http.Request) {
	returnString := "The application is online"
	response, _ := json.Marshal(returnString)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

//Authentication middleware function
func authMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		urlListToNotCheck := []string{"/login", "/refresh", "/bootstrap"}
		for _, uri := range urlListToNotCheck {
			if r.RequestURI == uri {
				h.ServeHTTP(w, r)
				return
			}
		}

		//Checks to see if a valid jwt token is present, if not will send a response back to client
		//and return a bool if request was successfully authroized
		isAuthenticated := validateToken(w, r)
		//Once passed, request is sent along to our route function
		if isAuthenticated {
			h.ServeHTTP(w, r)
		}
	})
}
