package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	db "orchestrallm/db"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
)

func TestApiKeyDatabaseIDsAndDelete(t *testing.T) {
	db.InitDB(filepath.Join(t.TempDir(), "keys.db"))
	aks = []ApiKey{}
	t.Cleanup(func() { aks = []ApiKey{} })
	t.Setenv("ORCHESTRALLM_TEST_KEY", "first")
	t.Setenv("ORCHESTRALLM_TEST_SECOND", "second")

	firstID, err := db.SaveDataAK("encrypted-first", "ORCHESTRALLM_TEST_KEY")
	if err != nil {
		t.Fatal(err)
	}
	secondID, err := db.SaveDataAK("encrypted-second", "ORCHESTRALLM_TEST_SECOND")
	if err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/ak", getApiKey).Methods("GET")
	router.HandleFunc("/ak/{id}", deleteApiKey).Methods("DELETE")
	request := func(method, path string, want int) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(method, path, nil))
		if w.Code != want {
			t.Fatalf("%s %s: got %d, want %d: %s", method, path, w.Code, want, w.Body.String())
		}
		return w
	}

	// Delete must work after restart, with an empty slice and no preceding GET.
	request(http.MethodDelete, "/ak/"+firstID, http.StatusNoContent)
	if _, exists := os.LookupEnv("ORCHESTRALLM_TEST_KEY"); exists {
		t.Fatal("deleted environment variable still exists")
	}
	keys, err := db.GetDataAK()
	if err != nil || len(keys) != 1 || keys[0].ID != secondID {
		t.Fatalf("wrong database contents after deletion: %v, %v", keys, err)
	}
	request(http.MethodDelete, "/ak/"+firstID, http.StatusNotFound)

	thirdID, err := db.SaveDataAK("encrypted-third", "ORCHESTRALLM_TEST_THIRD")
	if err != nil {
		t.Fatal(err)
	}
	if thirdID == firstID || thirdID == secondID {
		t.Fatal("database ID was reused")
	}
	aks = []ApiKey{}
	w := request(http.MethodGet, "/ak", http.StatusOK)
	var loaded []ApiKey
	if err := json.Unmarshal(w.Body.Bytes(), &loaded); err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 || loaded[0].ID != secondID || loaded[1].ID != thirdID {
		t.Fatalf("GET did not return all database IDs: %v", loaded)
	}
	request(http.MethodDelete, "/ak/"+secondID, http.StatusNoContent)
	request(http.MethodDelete, "/ak/"+thirdID, http.StatusNoContent)
	w = request(http.MethodGet, "/ak", http.StatusOK)
	if w.Body.String() != "[]\n" {
		t.Fatalf("empty database returned a fake record: %s", w.Body.String())
	}
}
