package main

import (
	"encoding/json"
	"fmt"
	"github.com/mariusmagureanu/burlan/src/pkg/log"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mariusmagureanu/burlan/src/pkg/entities"
)

func reqWrapper(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Add("Access-Control-Request-Method", "GET")
		w.Header().Add("Access-Control-Expose-Headers", "X-Jwt")

		log.Debug(r.Method, "", r.URL.RequestURI())
		fn(w,r)
	}
}
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

func login(w http.ResponseWriter, r *http.Request) {
	u := mux.Vars(r)["name"]
	var user entities.User

	err := db.Users().GetByName(&user, u)

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	token, err := jwtWrapper.GenerateToken(user.UID, user.Name, user.Email)

	if err != nil {
		log.Warning(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	w.Header().Set("X-JWT", token)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []entities.User

	err := db.Users().GetAll(&users)

	if err != nil {
		log.Error("db", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(users)

	if err != nil {
		log.Error("json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(out)
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
