package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	db "orchestrallm/db"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	fmt.Println("we've started...")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Opening db...")
	db.InitDB(filepath.Join(homeDir, "orchestrallm.db"))
	fmt.Println("DB Ready...")

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/index.html")
	}).Methods("GET")
	r.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		currentUser, err := user.Current()
		if err != nil {
			http.Error(w, "cannot get username", http.StatusInternalServerError)
			return
		}
		page, err := template.ParseFiles("html/main.html")
		if err != nil {
			http.Error(w, "cannot load page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := page.Execute(w, currentUser.Username); err != nil {
			fmt.Println(err)
		}
	})
	r.HandleFunc("/ak", getApiKey).Methods("GET")
	r.HandleFunc("/ak", postApiKey).Methods("POST")
	r.HandleFunc("/ak/{id}", putApiKey).Methods("PUT")
	r.HandleFunc("/ak/{id}", deleteApiKey).Methods("DELETE")
	r.HandleFunc("/chat", chatHandler).Methods("POST")

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
