package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var signals = []ModelSignal{}
var signalMu sync.RWMutex
var signalExecMu sync.Mutex
var sigCount int

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
	runSignals(w, r, []ModelSignal{signal})
}

func runSignals(w http.ResponseWriter, r *http.Request, incoming []ModelSignal) {
	signalExecMu.Lock()
	defer signalExecMu.Unlock()

	signalMu.Lock()
	start := len(signals)
	signals = append(signals, incoming...)
	pending := append([]ModelSignal{}, signals...)
	signalMu.Unlock()

	for i := sigCount; i < len(pending); i++ {
		prompt := fmt.Sprintf("Give a raw answer without formatation. User's prompt: %s", pending[i].Prompt)
		if i > 0 {
			prompt += "\nPrevious model's answer:\n" + pending[i-1].Response
		}
		resp := requestChat(w, r, pending[i].Model, prompt)
		if resp == nil {
			return
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			resp.Body.Close()
			return
		}
		var result struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		err := json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if err != nil || len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
			http.Error(w, "model returned no answer", http.StatusBadGateway)
			return
		}

		if

		pending[i].Response = result.Choices[0].Message.Content
		signalMu.Lock()
		signals[i].Response = pending[i].Response
		signalMu.Unlock()
		sigCount++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pending[start:])
}
