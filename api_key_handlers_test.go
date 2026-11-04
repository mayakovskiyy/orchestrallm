package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestApiKeyLifecycle(t *testing.T) {
	aks = []ApiKey{}
	nextApiKeyID = 0
	t.Cleanup(func() { aks = []ApiKey{}; nextApiKeyID = 0 })
	t.Setenv("ORCHESTRALLM_TEST_KEY", "")
	t.Setenv("ORCHESTRALLM_TEST_RENAMED", "")
	router := mux.NewRouter()
	router.HandleFunc("/ak", postApiKey).Methods("POST")
	router.HandleFunc("/ak/{id}", putApiKey).Methods("PUT")
	router.HandleFunc("/ak/{id}", deleteApiKey).Methods("DELETE")
	request := func(method, path, body string, want int) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
		if w.Code != want {
			t.Fatalf("%s %s: got %d, want %d: %s", method, path, w.Code, want, w.Body.String())
		}
		return w
	}
	request(http.MethodPost, "/ak", `{"name":"orchestrallm_test_key","key":"first"}`, http.StatusCreated)
	request(http.MethodPut, "/ak/1", `{`, http.StatusBadRequest)
	request(http.MethodPut, "/ak/1", `{"name":"","key":"bad"}`, http.StatusBadRequest)
	if os.Getenv("ORCHESTRALLM_TEST_KEY") != "first" {
		t.Fatal("invalid update changed the environment")
	}
	w := request(http.MethodPut, "/ak/1", `{"id":"ignored","name":"orchestrallm_test_renamed","key":"second"}`, http.StatusOK)
	var updated ApiKey
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil || updated.ID != "1" {
		t.Fatalf("update did not preserve ID: %s", w.Body.String())
	}
	if _, exists := os.LookupEnv("ORCHESTRALLM_TEST_KEY"); exists {
		t.Fatal("old environment variable still exists")
	}
	if os.Getenv("ORCHESTRALLM_TEST_RENAMED") != "second" {
		t.Fatal("updated environment variable is missing")
	}
	request(http.MethodDelete, "/ak/1", "", http.StatusNoContent)
	if _, exists := os.LookupEnv("ORCHESTRALLM_TEST_RENAMED"); exists {
		t.Fatal("deleted environment variable still exists")
	}
	request(http.MethodDelete, "/ak/1", "", http.StatusNotFound)
	request(http.MethodPut, "/ak/1", `{}`, http.StatusNotFound)
	request(http.MethodPost, "/ak", `{"name":"orchestrallm_test_key","key":"third"}`, http.StatusCreated)
	if len(aks) != 1 || aks[0].ID != "2" {
		t.Fatal("deleted ID was reused")
	}
}
