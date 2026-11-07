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
var nextApiKeyID uint64

func getApiKey(w http.ResponseWriter, r *http.Request) {
	var ak ApiKey
	aksMu.RLock()
	defer aksMu.RUnlock()

	if len(aks) == 0 {
		key, provider := db.GetDataAK()
		ak.Key = key
		ak.Name = provider
		aks = append(aks, ak)
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
	nextApiKeyID++
	ak.ID = fmt.Sprintf("%d", nextApiKeyID)
	aks = append(aks, ak)

	ck, err := enc.Encrypt(ak.Key, gk())
	if err != nil {
		fmt.Println(err)
	}

	if ck == "" {
		return
	}

	db.SaveDataAK(ck, strings.ToUpper(ak.Name))

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
	for i, ak := range aks {
		if ak.ID != mux.Vars(r)["id"] {
			continue
		}
		os.Unsetenv(strings.ToUpper(ak.Name))
		aks = append(aks[:i], aks[i+1:]...)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, "API key not found", http.StatusNotFound)
}
