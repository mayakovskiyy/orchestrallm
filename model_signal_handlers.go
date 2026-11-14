package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

var signals = []ModelSignal{}
var signalMu sync.RWMutex

func getSignal(w http.ResponseWriter, r *http.Request) {
	signalMu.RLock()
	defer signalMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signals)
}

func postSignal(w http.ResponseWriter, r *http.Request) {
	var signal ModelSignal
	if err := json.NewDecoder(r.Body).Decode(&signal); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	signalMu.Lock()

	defer signalMu.Unlock()

	signals = append(signals, signal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signal)
}
