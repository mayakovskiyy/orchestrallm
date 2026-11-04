package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type FirstResp struct {
	Response string `json:"response"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FirstResp{Response: "Welcome to the OrchestraLLM API! Visit http://localhost:8000/help to learn more about OrchestraLLM. Made by maya."})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
