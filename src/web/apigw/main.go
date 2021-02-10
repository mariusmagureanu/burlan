package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mariusmagureanu/burlan/src/pkg/auth"
	"github.com/mariusmagureanu/burlan/src/pkg/dao"
	"github.com/mariusmagureanu/burlan/src/pkg/log"
)

const (
	apiNameSpace = "/api/v1/"

	rsaKey = "MIIBOQIBAAJAdr72N291gUt/nticv6z46YzCglpYLruGXyyvV/mTBXqWKojw5XlY\nZAd2RmHFLrr1fgGPo2uDkWC03O+/pENAdQIDAQABAkBsZXu7NQrN4U55gYDtVAfQosa4WaJf3p0V6mOR6mh0Oaj4DdgffS/UoaeIuCIEVIbXUN7ndXUk0aeD/XNZU1HhAiEAvB0gEvHQoS21pEwTlYohsEFz1cLHRujChxp4cHCoAq0CIQChmUu3eObdWUYdtWDaVg+qSPXJE1xZVhmD4zbqUlc16QIgMvf9QcTNT26QIbUPNVxY5mXFmeyNi/PzCSIt8eFEVH0CIBUF6nHOCsrVKGgJBrag55zRrRghqqv8pYkg8C3/1FSxAiEAl7bipiY5nK+cJQ8XpJFT0QWRu1MufuYFf+MiVKf7zTU="
)

var (
	jwtWrapper        = auth.JwtWrapper{}
	db                = dao.DAO{}
	commandLine       = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	versionFlag       = commandLine.Bool("V", false, "Show version and exit")
	portFlag          = commandLine.Uint("port", uint(8080), "Port used by the http server")
	hostFlag          = commandLine.String("host", "localhost", "Host for the http server")
	brokersFlag       = commandLine.String("brokers", "localhost:9092", "Kafka addresses separate by comma, if multiple specified")
	logLevelFlag      = commandLine.String("log", "debug", "Available levels: debug,info,warning,error,quiet")
	jwtExpirationTime = commandLine.Int64("jwt-exp", 24, "JWT expiration time expressed in hours")

	version  = "N/A"
	revision = "N/A"
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

	log.InitNewLogger(os.Stdout, log.GetLogLevelID(*logLevelFlag))

	jwtWrapper.SecretKey = rsaKey
	jwtWrapper.ExpirationHours = *jwtExpirationTime
	hostname, err := os.Hostname()

	if err != nil {
		log.ErrorSync(err.Error())
		os.Exit(1)
	}

	jwtWrapper.Issuer = hostname

	err = db.Init("foo.sqlite")

	if err != nil {
		log.ErrorSync(err)
		os.Exit(1)
	}

	err = db.CreateTables()
	if err != nil {
		log.ErrorSync(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	hub := newHub()
	go hub.run()

	r.HandleFunc(apiNameSpace+"users", reqWrapper(getAllUsers)).Methods(http.MethodGet,http.MethodOptions)
	r.HandleFunc(apiNameSpace+"user", reqWrapper(createUser)).Methods(http.MethodPost)
	r.HandleFunc(apiNameSpace+"user/{id}", reqWrapper(getUserByID)).Methods(http.MethodGet)
	r.HandleFunc(apiNameSpace+"login/{name}", reqWrapper(login)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	brokers = strings.Split(*brokersFlag, ",")
	addr := fmt.Sprintf("%s:%d", *hostFlag, *portFlag)

	log.InfoSync(fmt.Sprintf("Started listening on: <%s>", addr))
	log.ErrorSync(http.ListenAndServe(addr, r))
}
