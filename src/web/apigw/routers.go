package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mariusmagureanu/burlan/src/pkg/entities"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	mime := r.Header.Get("Content-Type")

	if mime != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var user entities.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Users().Insert(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Location", fmt.Sprintf("/user/%d", user.ID))
	w.WriteHeader(http.StatusCreated)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	i, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user entities.User

	err = db.Users().GetByID(&user, uint(i))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(out)

}

func addFriendToUser(w http.ResponseWriter, r *http.Request) {
}

func removeFriendFromUser(w http.ResponseWriter, r *http.Request) {
}
