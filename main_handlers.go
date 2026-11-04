package main

import (
	"encoding/json"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FirstResp{Response: "Welcome to the OrchestraLLM API! Visit http://localhost:8000/help to learn more about OrchestraLLM. Made by maya."})
}
