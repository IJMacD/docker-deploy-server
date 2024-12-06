package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path"
)

func serveRevision(w http.ResponseWriter, r	*http.Request, f string, rn string) {
	if f == "" {
		f = defaultFleet
	}

	if rn == "" {
		rn = defaultRevision
	}

	filename := path.Join("fleets", f, rn+".yml")
	file, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/yaml")

	hash := md5.Sum([]byte(file))
	etag := `"`+hex.EncodeToString(hash[:])+`"`

	rEtag := r.Header.Get("If-None-Match")

	if rEtag == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("ETag", etag)

	w.Write(file)
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
