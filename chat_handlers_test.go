package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatRejectsInvalidRequests(t *testing.T) {
	for _, body := range []string{
		`{`,
		`{"model":"openrouter/free","prompt":" "}`,
		`{"model":"paid/model","prompt":"hello"}`,
	} {
		w := httptest.NewRecorder()
		chatHandler(w, httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(body)))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("got %d for %s", w.Code, body)
		}
	}
}
