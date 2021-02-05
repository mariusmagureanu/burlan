package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mariusmagureanu/burlan/src/pkg/dao"
	"github.com/mariusmagureanu/burlan/src/pkg/log"
)

const apiNameSpace = "/api/v1/"

var (
	db          = dao.DAO{}
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	versionFlag = commandLine.Bool("V", false, "Show version and exit")
	portFlag    = commandLine.Uint("port", uint(8080), "Port used by the http server")
	hostFlag    = commandLine.String("host", "localhost", "Host for the http server")
	brokersFlag = commandLine.String("brokers", "localhost:9092", "Kafka addresses separate by comma, if multiple specified")
	logLevelFlag=commandLine.String("log", "debug", "Available levels: debug,info,warning,error,quiet")
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
		log.ErrorSync(err)
		os.Exit(1)
	}

	if *versionFlag {
		fmt.Println("Version:  " + version)
		fmt.Println("Revision: " + revision)
		os.Exit(0)
	}

	log.InitNewLogger(os.Stdout, log.ErrorLevel)
	log.SetLogLevel(log.GetLogLevelID(*logLevelFlag))

	err := db.Init("foo.sqlite")

	if err != nil {
		log.ErrorSync(err)
		os.Exit(1)
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

	brokers = strings.Split(*brokersFlag, ",")
	addr := fmt.Sprintf("%s:%d", *hostFlag, *portFlag)

	log.InfoSync(fmt.Sprintf("Started listening on: <%s>",addr))
	log.ErrorSync(http.ListenAndServe(addr, r))
}
