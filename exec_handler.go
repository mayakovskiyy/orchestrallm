package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func postExec(w http.ResponseWriter, r *http.Request) {
	var signal ExecSignal
	if err := json.NewDecoder(r.Body).Decode(&signal); err != nil || strings.TrimSpace(signal.Prompt) == "" || len(signal.Models) == 0 {
		http.Error(w, "enter a prompt and a list of models", http.StatusBadRequest)
		return
	}
	models, err := json.Marshal(signal.Models)
	if err != nil {
		http.Error(w, "cannot encode models", http.StatusInternalServerError)
		return
	}
	testJson := `[{"id": "1", "model": "example/example", "role": "Programmer", "prompt": "Task"}, {"id": "2", "model": "other/model", "role": "Code-review", "prompt": "Review task"}]`
	prompt := fmt.Sprintf("First of all: DON'T WRITE ANY TEXT OUT FROM THE JSON REQUEST YOU'LL WRITE (DO NOT. WRITE THE ONLY JSON REQUEST.). Answer like a professional Team Lead.' You have a list of the models and you should: A) Pick a codder-model B) Pick a code-review model. Write full technical requirements and tasks for other models, recheck your answer if you're wrong. Avoid free Google models. And also your answer should be based on JSON value signal and contain: ID (ID IS A FIELD. 1, 2, 3 etc.), Model (indicated in an array, also choose the name of the model from OpenRouter. Models SHOULDN'T REPEAT. One model tooks one role. You shouldn't make more than a one Programmer or a Code-reviewer (and other roles). Code-reviewer should fix the code and return an already fixed code. If user doesn't indicate about their own roles and agents - don't use them and use only Programmer and Code-reviewer.), Role (If the model is a programmer write 'Programmer'; 'Code-review' etc. Programmer and Code-reviewer are required.), Prompt (your direct actions and what should a model do. Indicate their role in the text.). Just give a raw JSON as a text. DO NOT formate it as code. User prompt: %s The array of the available models (The field in JSON shouldn't be an array. Just an ordinary string-field): %s Return one JSON array of signals, without Markdown, in this style: %s", signal.Prompt, models, testJson)
	response := requestChat(w, r, signal.Model, prompt)
	if response == nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.WriteHeader(response.StatusCode)
		io.Copy(w, response.Body)
		return
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil || len(result.Choices) == 0 {
		http.Error(w, "model returned no answer", http.StatusBadGateway)
		return
	}
	parsed, err := parseExecSignals(result.Choices[0].Message.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	runSignals(w, r, parsed)
}

func parseExecSignals(content string) ([]ModelSignal, error) {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	var parsed []ModelSignal
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("invalid signals JSON: %w", err)
	}
	if parsed == nil {
		return nil, fmt.Errorf("expected a JSON array of signals")
	}
	for _, signal := range parsed {
		if signal.ID == "" || signal.Model == "" || signal.Role == "" || signal.Prompt == "" {
			return nil, fmt.Errorf("signal must contain id, model, role and prompt")
		}
	}
	return parsed, nil
}
