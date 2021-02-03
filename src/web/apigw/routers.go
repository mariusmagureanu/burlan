package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func foo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	page := vars["page"]

	_, _ = fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
}
