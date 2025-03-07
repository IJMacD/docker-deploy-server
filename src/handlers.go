package main

import (
	"fmt"
	"net/http"
	"path"
	"regexp"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "index", struct {
		Fleets   []Fleet
		Machines []*Machine
	}{Fleets: fleets, Machines: machines})
}

func revisionHandler(w http.ResponseWriter, r *http.Request, fleet string, revision string) {
	if r.Method == "POST" {
		saveHandler(w, r, fleet, revision)
		return
	}

	// Handle special case of creating new revision
	if revision == getNextRevisionName(fleet) {
		editHandler(w, r, fleet, revision)
		return
	}

	viewHandler(w, r, fleet, revision)
}

func viewHandler(w http.ResponseWriter, _ *http.Request, fleet string, revision string) {
	rev, err := getRevision(fleet, revision)

	if err != nil {
		http.Error(w, "Revision not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "view", rev)
}

func editHandler(w http.ResponseWriter, r *http.Request, fleet string, revision string) {
	name := getNextRevisionName(fleet)

	if name != revision {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	rev := Revision{Name: name, FleetName: fleet}

	from := r.URL.Query().Get("from")
	if from != "" {
		v, err := getRevision(rev.FleetName, from)
		if err == nil {
			rev.Body = v.Body
		}
	}

	renderTemplate(w, "edit", rev)
}

func saveHandler(w http.ResponseWriter, r *http.Request, fleet string, revision string) {
	name := getNextRevisionName(fleet)

	if name != revision {
		http.Error(w, "Revision does not match expected revision", http.StatusBadRequest)
		return
	}

	body := r.FormValue("body")
	p := Revision{Name: name, FleetName: fleet, Body: []byte(body)}

	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/revisions/"+fleet+"/"+name, http.StatusFound)
}

func fleetViewHandler(w http.ResponseWriter, r *http.Request, fleet string) {
	f, ok := getFleet(fleet)

	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "fleet-view", f)
}

func fleetCreateHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	_, err := regexp.Match(`^[a-zA-Z0-9-]+$`, []byte(name))

	_, ok := getFleet(name)

	if err == nil && !ok {
		err = makeFleet(name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func machineHandler(w http.ResponseWriter, r *http.Request, serial string) {
	fleet := r.FormValue("fleet")

	m, ok := getMachine(serial)

	if ok {
		m.fleetName = fleet
		m.RevisionName = ""
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	m := validMachineApiPath.FindStringSubmatch(r.URL.Path)
	if m != nil {
		apiHandlerMachine(w, r, m[1])
		return
	}

	m = validFleetApiPath.FindStringSubmatch(r.URL.Path)
	if m != nil {
		apiHandlerFleet(w, r, m[1])
		return
	}

	http.NotFound(w, r)
}

func apiHandlerMachine(w http.ResponseWriter, r *http.Request, m string) {
	n, ok := getMachine(m)
	if !ok {
		n = &Machine{SerialNumber: m}
		machines = append(machines, n)
	}

	f, ok := getFleet(n.fleetName)
	if !ok {
		f = &Fleet{}
	}

	// Record latest revision served to machine
	n.RevisionName = f.CurrentRevisionName()
	n.LastSync = time.Now().Format(time.RFC3339)

	http.Redirect(w, r, fmt.Sprintf("/api/v1/fleets/%s/docker-compose.yml", n.FleetName()), http.StatusFound)
}

func apiHandlerFleet(w http.ResponseWriter, r *http.Request, n string) {
	f, ok := getFleet(n)
	if !ok {
		f = &Fleet{}
	}

	serveRevision(w, r, n, f.CurrentRevisionName())
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	m := validStaticPath.FindStringSubmatch(r.URL.Path)

	if m == nil {
		http.NotFound(w, r)
		return
	}

	filename := path.Join("static", m[1])
	http.ServeFile(w, r, filename)
}
