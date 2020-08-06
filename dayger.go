package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bdkiran/dayger/api"
)

func main() {
	log.Println("Welcome to dayger")
	handleRouter()
}

func handleRouter() {
	router := api.NewRouter()

	port := ":8080"
	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
