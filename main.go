package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Team-Udaan/udaan16-gslr-poll-api/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	http.ListenAndServe(":8081", router)
}