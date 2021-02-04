package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/mariusmagureanu/burlan/src/pkg/dao"
)

const apiNameSpace = "/api/v1/"

var (
	db          = dao.DAO{}
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	versionFlag = commandLine.Bool("V", false, "Show version and exit")
	version     = "N/A"
	revision    = "N/A"
)

func main() {
	commandLine.Usage = func() {
		fmt.Fprint(os.Stdout, "Usage of the api-gw:\n")
		commandLine.PrintDefaults()
		os.Exit(0)
	}

	if err := commandLine.Parse(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}

	if *versionFlag {
		fmt.Println("Version:  " + version)
		fmt.Println("Revision: " + revision)
		os.Exit(0)
	}

	err := db.Init("foo.sqlite")

	if err != nil {
		log.Fatalln(err)
	}

	db.CreateTables()

	r := mux.NewRouter()
	hub := newHub()
	go hub.run()

	r.HandleFunc(apiNameSpace+"user", createUser).Methods(http.MethodPost)
	r.HandleFunc(apiNameSpace+"user/{id}", getUserByID).Methods(http.MethodGet)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	_ = http.ListenAndServe(":8080", r)
}
