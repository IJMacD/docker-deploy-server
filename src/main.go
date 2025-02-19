package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var fleets []Fleet
var machines = make([]*Machine, 0)

var machineLedger *os.File

var defaultFleet = "default"
var defaultRevision = "r0"

var templates = template.Must(template.ParseFiles(
	"tmpl/index.html",
	"tmpl/edit.html",
	"tmpl/view.html",
	"tmpl/fleet-view.html",
))
var validRevisionsPath = regexp.MustCompile(`^/revisions/([a-zA-Z0-9-]+)/(r\d+)$`)
var validFleetsPath = regexp.MustCompile(`^/fleets/([a-zA-Z0-9-]+)$`)
var validMachinesPath = regexp.MustCompile(`^/machines/([a-zA-Z0-9-]+)$`)
var validApiPath = regexp.MustCompile("^/api/machines/([a-zA-Z0-9-]+)/docker-compose.yml$")
var validStaticPath = regexp.MustCompile("^/static/([a-zA-Z0-9.-]+)$")

func main() {
	fs, err := makeFleets()
	if err != nil {
		panic("Unable to make fleets")
	}
	fleets = fs

	if _, ok := getFleet(defaultFleet); !ok {
		err = makeFleet(defaultFleet)

		if err != nil {
			panic("Unable to create default fleet")
		}
	}

	machineLedger, err = os.OpenFile("machines.csv", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0o600)
	if err != nil {
		panic("Unable to open machine ledger")
	}

	i, err := machineLedger.Stat()
	if err != nil {
		panic("Unable to determine ledger size")
	}
	if i.Size() == 0 {
		n, err := machineLedger.WriteString("machine,fleet,revision,lastSync\n")
		if err != nil {
			panic("Unable to write headers to machine ledger. " + err.Error())
		}
		fmt.Printf("Written %d bytes to %s\n", n, "machines.csv")
	}

	http.HandleFunc("/api/machines/", apiHandler)

	http.HandleFunc("/machines/", makeMachineHandler(machineHandler))

	http.HandleFunc("/revisions/", makeRevisionHandler(revisionHandler))

	http.HandleFunc("/fleets/", makeFleetHandler(fleetViewHandler))
	http.HandleFunc("/fleets", fleetCreateHandler)

	http.HandleFunc("/static/", staticHandler)

	http.HandleFunc("/", indexHandler)

	fmt.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", loggingHandler(http.DefaultServeMux)))
}
