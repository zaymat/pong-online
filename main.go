package main

import (
	"fmt"
	"net/http"

	"./handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Routing part
	r := mux.NewRouter()
	r.HandleFunc("/", home.HomeHandler)

	http.Handle("/", r)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", r)
}
