package main

import (
	"encoding/json"
	"fmt"
	// "github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"net/http"
	db "orchestrallm/db"
	enc "orchestrallm/encryptions"
	"os"
	"strings"
	"sync"
)

var aks = []ApiKey{}
var aksMu sync.RWMutex

func loadApiKeys() error {
	keys, err := db.GetDataAK()
	if err != nil {
		return err
	}
	aks = []ApiKey{}
	for _, ak := range keys {
		aks = append(aks, ApiKey{ID: ak.ID, Key: ak.Key, Name: ak.Name})
	}
	return nil
}

func getApiKey(w http.ResponseWriter, r *http.Request) {
	aksMu.Lock()
	defer aksMu.Unlock()
	if err := loadApiKeys(); err != nil {
		fmt.Println(err)
		http.Error(w, "database read failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aks)
}

func postApiKey(w http.ResponseWriter, r *http.Request) {
	var ak ApiKey
	if err := json.NewDecoder(r.Body).Decode(&ak); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	aksMu.Lock()
	defer aksMu.Unlock()
	if err := os.Setenv(strings.ToUpper(ak.Name), ak.Key); err != nil {
		http.Error(w, "invalid environment variable name or value", http.StatusBadRequest)
		return
	}
	ck, err := enc.Encrypt(ak.Key, gk())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "encryption failed", http.StatusInternalServerError)
		return
	}

	ak.ID, err = db.SaveDataAK(ck, strings.ToUpper(ak.Name))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "database write failed", http.StatusInternalServerError)
		return
	}
	aks = append(aks, ak)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ak)

}

func putApiKey(w http.ResponseWriter, r *http.Request) {
	var ak ApiKey
	if err := json.NewDecoder(r.Body).Decode(&ak); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	aksMu.Lock()
	defer aksMu.Unlock()
	if err := loadApiKeys(); err != nil {
		fmt.Println(err)
		http.Error(w, "database read failed", http.StatusInternalServerError)
		return
	}
	for i, existing := range aks {
		if existing.ID != mux.Vars(r)["id"] {
			continue
		}
		if err := os.Setenv(strings.ToUpper(ak.Name), ak.Key); err != nil {
			http.Error(w, "invalid environment variable name or value", http.StatusBadRequest)
			return
		}
		if strings.ToUpper(existing.Name) != strings.ToUpper(ak.Name) {
			os.Unsetenv(strings.ToUpper(existing.Name))
		}
		ak.ID = existing.ID
		aks[i] = ak

		ck, err := enc.Encrypt(ak.Key, gk())
		if err != nil {
			fmt.Println(err)
			http.Error(w, "encryption failed", http.StatusInternalServerError)
			return
		}
		db.UpdateDataAK(ck, strings.ToUpper(ak.Name), existing.ID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ak)
		return
	}
	http.Error(w, "API key not found", http.StatusNotFound)
}

func deleteApiKey(w http.ResponseWriter, r *http.Request) {
	aksMu.Lock()
	defer aksMu.Unlock()
	if err := loadApiKeys(); err != nil {
		fmt.Println(err)
		http.Error(w, "database read failed", http.StatusInternalServerError)
		return
	}
	for i, ak := range aks {
		if ak.ID != mux.Vars(r)["id"] {
			continue
		}
		if err := db.DeleteDataAK(ak.ID); err != nil {
			fmt.Println(err)
			http.Error(w, "database delete failed", http.StatusInternalServerError)
			return
		}
		os.Unsetenv(strings.ToUpper(ak.Name))
		aks = append(aks[:i], aks[i+1:]...)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, "API key not found", http.StatusNotFound)
}
