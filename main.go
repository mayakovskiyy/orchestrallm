package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/ak", getApiKey).Methods("GET")
	r.HandleFunc("/ak", postApiKey).Methods("POST")
	r.HandleFunc("/ak/{id}", putApiKey).Methods("PUT")
	r.HandleFunc("/ak/{id}", deleteApiKey).Methods("DELETE")

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
