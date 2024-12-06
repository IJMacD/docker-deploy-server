package main

import (
	"os"
	"path"
)

func getFleetNames() ([]string, error) {
	l, err := os.ReadDir("fleets")
	if err != nil {
		return nil, err
	}

	var r []string

	for _, f := range l {
		if f.Type().IsDir() {
			r = append(r, f.Name())
		}
	}

	return r, nil
}

func makeFleets() ([]Fleet, error) {
	fleets := make([]Fleet, 0)

	fleetNames, err := getFleetNames()

	if err != nil {
		return nil, err
	}

	for _, f := range fleetNames {
		fleets = append(fleets, Fleet{Name: f})
	}

	return fleets, nil
}

func getFleet(name string) (*Fleet, bool) {
	if name == "" {
		name = defaultFleet
	}

	for i, f := range fleets {
		if f.Name == name || (f.Name == "" && name == defaultFleet) {
			return &fleets[i], true
		}
	}

	return nil, false
}

func makeFleet(name string) error {
	err := os.Mkdir(path.Join("fleets", name), 0o700)

	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join("fleets", name, defaultRevision+".yml"), []byte("services:\n"), 0o600)

	if err != nil {
		return err
	}

	fleets = append(fleets, Fleet{Name: name})

	return nil
}
