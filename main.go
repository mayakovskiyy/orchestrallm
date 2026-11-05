package main

import (
	//	"encoding/json"
	"fmt"
	"net/http"
	db "orchestrallm/db"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	db.InitDB(homeDir)

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/index.html")
	}).Methods("GET")
	r.HandleFunc("/ak", getApiKey).Methods("GET")
	r.HandleFunc("/ak", postApiKey).Methods("POST")
	r.HandleFunc("/ak/{id}", putApiKey).Methods("PUT")
	r.HandleFunc("/ak/{id}", deleteApiKey).Methods("DELETE")

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
