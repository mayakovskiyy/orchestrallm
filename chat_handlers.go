package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	db "orchestrallm/db"
	enc "orchestrallm/encryptions"
	"strings"
	"time"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || strings.TrimSpace(input.Prompt) == "" {
		http.Error(w, "enter a prompt", http.StatusBadRequest)
		return
	}
	response := requestChat(w, r, input.Model, input.Prompt)
	if response == nil {
		return
	}
	defer response.Body.Close()
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

func requestChat(w http.ResponseWriter, r *http.Request, model string, prompt string) *http.Response {
	if model != "openrouter/free" && !strings.HasSuffix(model, ":free") {
		http.Error(w, "select a free model", http.StatusBadRequest)
		return nil
	}

	keys, err := db.GetDataAK()
	if err != nil {
		http.Error(w, "cannot load API key", http.StatusInternalServerError)
		return nil
	}
	var key string
	for _, ak := range keys {
		name := strings.ToUpper(strings.TrimSpace(ak.Name))
		if name == "OPENROUTER" || name == "OPENROUTER_API_KEY" {
			key, err = enc.Decrypt(ak.Key, gk())
			if err != nil {
				http.Error(w, "cannot decrypt OpenRouter key", http.StatusInternalServerError)
				return nil
			}
			break
		}
	}
	if key == "" {
		http.Error(w, "add an OpenRouter key in the dashboard first (name: OPENROUTER)", http.StatusBadRequest)
		return nil
	}

	body, err := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	})
	if err != nil {
		http.Error(w, "cannot encode request", http.StatusInternalServerError)
		return nil
	}
	req, err := http.NewRequestWithContext(r.Context(), "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		http.Error(w, "cannot create request", http.StatusInternalServerError)
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 120 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, "OpenRouter request failed", http.StatusBadGateway)
		return nil
	}
	return response
}
