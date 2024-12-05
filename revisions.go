package main

import (
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func getRevisionNames(f string) ([]string, error) {
	if f == "" {
		f = defaultFleet
	}
	
	l, err := os.ReadDir(path.Join("fleets", f))
	if err != nil {
		return nil, err
	}

	r := make([]string, 0)

	x := regexp.MustCompile(`^r\d+\.yml$`)

	for _, f := range l {
		if x.Match([]byte(f.Name())) && f.Type().IsRegular() {
			r = append(r, strings.Replace(f.Name(), ".yml", "", 1))
		}
	}

	slices.Reverse(r)

	return r, nil
}

func getLatestRevisionName(fleet string) string {
	rs, err := getRevisionNames(fleet)

	if err != nil || len(rs) == 0 {
		return defaultRevision
	}

	return rs[0]
}

func getNextRevisionName(fleet string) string {
	l := getLatestRevisionName(fleet)

	n, err := strconv.Atoi(strings.TrimPrefix(l, "r"))

	if err != nil {
		return defaultRevision
	}

	return "r" + strconv.Itoa(int(n)+1)
}

func getRevision(f string, r string) (*Revision, error) {
	if f == "" {
		f = defaultFleet
	}

	if r == "" {
		r = defaultRevision
	}

	v := Revision{Name: r, FleetName: f}

	filename := path.Join("fleets", f, r+".yml")
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	v.Body = b

	return &v, nil
}