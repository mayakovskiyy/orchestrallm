package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseExecSignals(t *testing.T) {
	content := `[{"id":"1","model":"example/model:free","role":"Programmer","prompt":"Task"}]`
	for _, input := range []string{content, "```json\n" + content + "\n```", "```\n" + content + "\n```"} {
		parsed, err := parseExecSignals(input)
		if err != nil || len(parsed) != 1 || parsed[0].Prompt != "Task" {
			t.Fatalf("cannot parse signal: %v", err)
		}
	}
	for _, input := range []string{"", "null", `{}`, `[{}]`, `[{"id":1}]`} {
		if _, err := parseExecSignals(input); err == nil {
			t.Fatalf("accepted invalid signals: %s", input)
		}
	}
}

func TestExecRejectsInvalidInput(t *testing.T) {
	for _, body := range []string{
		`{`,
		`{"prompt":" ","models":["openrouter/free"]}`,
		`{"prompt":"Task","models":[]}`,
		`{"prompt":"Task","models":["openrouter/free"],"model":"paid/model"}`,
	} {
		w := httptest.NewRecorder()
		postExec(w, httptest.NewRequest(http.MethodPost, "/api/v1/orchestra/exec", strings.NewReader(body)))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("got %d for %s", w.Code, body)
		}
	}
}
