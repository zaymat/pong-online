package home

import (
	"fmt"
	"net/http"
)

// HomeHandler : handler for home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, welcome to the server !")
}
