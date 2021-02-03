package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", foo)

	_ = http.ListenAndServe(":8080", r)
}
