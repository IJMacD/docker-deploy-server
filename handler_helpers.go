package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

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

func makeFleetHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validFleetsPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

func makeRevisionHandler(fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validRevisionsPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1], m[2])
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

func renderTemplate(w http.ResponseWriter, tmpl string, p any) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
