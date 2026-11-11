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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aks)
}

func postSignal(w http.ResponseWriter, r *http.Request) {
	var signal ModelSignal
	signalMu.Lock()

	defer signalMu.Unlock()

	signals = append(signals, signal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signal)
}
