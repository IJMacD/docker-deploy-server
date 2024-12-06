package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var fleets []Fleet
var machines = make([]*Machine, 0)

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
