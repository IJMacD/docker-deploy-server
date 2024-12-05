package main

import (
	"net/http"
	"os"
	"path"
	"regexp"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", struct {
		Fleets  	[]Fleet
		Machines  	[]*Machine
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
	
	http.Redirect(w, r, "/revisions/" +fleet+ "/"+name, http.StatusFound)
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
		err = os.Mkdir(path.Join("fleets", name), 0o700)

		if err == nil {
			fleets = append(fleets, Fleet{Name: name})

			err = os.WriteFile(path.Join("fleets", name, defaultRevision+".yml"), []byte("services:\n"), 0o600)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func machineHandler(w http.ResponseWriter, r *http.Request, serial string) {
	fleet := r.FormValue("fleet")

	m, ok := getMachine(serial)

	f, ok2 := getFleet(fleet)

	if ok && ok2 {
		m.FleetName = fleet
		m.RevisionName = f.CurrentRevisionName()
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	m := validApiPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	n, ok := getMachine(m[1])
	if !ok {
		n = &Machine{SerialNumber: m[1]}
		machines = append(machines, n)
	}

	f, ok := getFleet(n.FleetName)
	if !ok {
		f = &Fleet{}
	}

	// Record latest revision served to machine
	n.RevisionName = f.CurrentRevisionName()

	serveRevision(w, r, n.FleetName, n.RevisionName)
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
