package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignalReadThenWrite(t *testing.T) {
	signals = []ModelSignal{}
	sigCount = 0
	t.Cleanup(func() { signals = []ModelSignal{}; sigCount = 0 })
	w := httptest.NewRecorder()
	getSignal(w, httptest.NewRequest(http.MethodGet, "/api/v1/orchestra/signal", nil))
	if w.Body.String() != "[]\n" {
		t.Fatalf("expected empty signals, got %s", w.Body.String())
	}
	if !signalMu.TryLock() {
		t.Fatal("GET did not release the read lock")
	}
	signalMu.Unlock()

	w = httptest.NewRecorder()
	body := `{"id":"1","model":"example/model","role":"Programmer","prompt":"Task"}`
	postSignal(w, httptest.NewRequest(http.MethodPost, "/api/v1/orchestra/signal", strings.NewReader(body)))
	// An unsupported model fails before an external request, keeping the signal pending.
	if w.Code != http.StatusBadRequest {
		t.Fatalf("POST returned %d: %s", w.Code, w.Body.String())
	}
	if sigCount != 0 {
		t.Fatal("failed signal was marked as processed")
	}
	w = httptest.NewRecorder()
	getSignal(w, httptest.NewRequest(http.MethodGet, "/api/v1/orchestra/signal", nil))
	var result []ModelSignal
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0] != signals[0] || result[0].Prompt != "Task" {
		t.Fatalf("GET did not return the posted signal: %s", w.Body.String())
	}
}
