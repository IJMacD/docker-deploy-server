package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Revision struct {
	Name  string
	Body  []byte
	Fleet string
}

type Machine struct {
	SerialNumber string
	Fleet        string
	Revision     string
}

var machines = make([]Machine, 0)

var defaultFleet = "default"
var defaultRevision = "r0"

var templates = template.Must(template.ParseFiles("tmpl/index.html", "tmpl/edit.html", "tmpl/view.html"))
var validRevisionsPath = regexp.MustCompile(`^/(edit|save|revisions)/(r\d+)$`)
var validMachinesPath = regexp.MustCompile(`^/machines/([a-zA-Z0-9]+)$`)
var validApiPath = regexp.MustCompile("^/api/machines/([a-zA-Z0-9]+)/docker-compose.yml$")
var validStaticPath = regexp.MustCompile("^/static/([a-zA-Z0-9.-]+)$")

func (r *Revision) save() error {
	filename := path.Join("fleets", r.Fleet, r.Name+".yml")
	return os.WriteFile(filename, r.Body, 0600)
}

func main() {
	http.HandleFunc("/api/machines/", apiHandler)

	http.HandleFunc("/machines/", makeMachineHandler(machineHandler))

	http.HandleFunc("/revisions/", makeRevisionHandler(viewHandler))
	http.HandleFunc("/edit/", makeRevisionHandler(editHandler))
	http.HandleFunc("/save/", makeRevisionHandler(saveHandler))

	http.HandleFunc("/static/", staticHandler)

	http.HandleFunc("/", indexHandler)

	fmt.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", loggingHandler(http.DefaultServeMux)))
}

type PageData struct {
	Machines  []Machine
	Revisions []string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	l, _ := getRevisionNames(defaultFleet)
	renderTemplate(w, "index", PageData{Machines: machines, Revisions: l})
}

func makeRevisionHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validRevisionsPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func makeMachineHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validMachinesPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

func loggingHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Printf("%s %s\n", r.Method, r.URL)
	}
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
		machines = append(machines, *n)
	}

	serveRevision(w, n.Fleet, n.Revision)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	m := validStaticPath.FindStringSubmatch(r.URL.Path)

	if m == nil {
		http.NotFound(w, r)
		return
	}

	filename := path.Join("static", m[1])
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "text/css")

	io.Copy(w, file)
}

func getMachine(serial string) (m *Machine, ok bool) {
	for i, m := range machines {
		if m.SerialNumber == serial {
			return &machines[i], true
		}
	}

	return nil, false
}

func getRevisionNames(f string) ([]string, error) {
	l, err := os.ReadDir(path.Join("fleets", f))
	if err != nil {
		return nil, err
	}

	var r []string

	x := regexp.MustCompile(`^r\d+\.yml$`)

	for _, f := range l {
		if x.Match([]byte(f.Name())) && f.Type().IsRegular() {
			r = append(r, strings.Replace(f.Name(), ".yml", "", 1))
		}
	}

	return r, nil
}

func getNextRevisionName(f string) string {
	rs, err := getRevisionNames(f)

	if err != nil {
		return defaultRevision
	}

	l := rs[len(rs)-1]

	n, err := strconv.Atoi(strings.TrimPrefix(l, "r"))

	if err != nil {
		return defaultRevision
	}

	return "r" + strconv.Itoa(int(n)+1)
}

func serveRevision(w http.ResponseWriter, f string, r string) {
	if f == "" {
		f = defaultFleet
	}

	if r == "" {
		r = defaultRevision
	}

	filename := path.Join("fleets", f, r+".yml")
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/yaml")

	io.Copy(w, file)
}

func getRevision(f string, r string) (Revision, error) {
	if f == "" {
		f = defaultFleet
	}

	if r == "" {
		r = defaultRevision
	}

	v := Revision{Name: r, Fleet: f}

	filename := path.Join("fleets", f, r+".yml")
	b, err := os.ReadFile(filename)
	if err != nil {
		return v, err
	}

	v.Body = b

	return v, nil
}

type ViewParams struct {
	Revision Revision
	NextName string
	Machines []Machine
}

func viewHandler(w http.ResponseWriter, r *http.Request, revision string) {
	rev, err := getRevision(defaultFleet, revision)

	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	name := getNextRevisionName(defaultFleet)

	ms := make([]Machine, 0)
	for _, m := range machines {
		if (m.Revision == revision || (m.Revision == "" && revision == defaultRevision)) && (m.Fleet == defaultFleet || m.Fleet == "") {
			ms = append(ms, m)
		}
	}

	renderTemplate(w, "view", ViewParams{Revision: rev, NextName: name, Machines: ms})
}

func editHandler(w http.ResponseWriter, r *http.Request, n string) {
	name := getNextRevisionName(defaultFleet)

	if name != n {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	rev := Revision{Name: name, Fleet: defaultFleet}

	from := r.URL.Query().Get("from")
	if from != "" {
		v, err := getRevision(rev.Fleet, from)
		if err == nil {
			rev.Body = v.Body
		}
	}

	renderTemplate(w, "edit", rev)
}

func saveHandler(w http.ResponseWriter, r *http.Request, n string) {
	name := getNextRevisionName(defaultFleet)

	body := r.FormValue("body")
	p := Revision{Name: name, Fleet: defaultFleet, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/revisions/"+name, http.StatusFound)
}

func machineHandler(w http.ResponseWriter, r *http.Request, serial string) {
	revision := r.FormValue("revision")

	m, ok := getMachine(serial)

	if ok {
		m.Revision = revision
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p any) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
