package main

import (
	//	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
