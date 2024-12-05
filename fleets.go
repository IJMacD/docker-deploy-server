package main

import "os"

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

	if (err != nil) {
		return nil, err
	}

	for _, f := range fleetNames {
		fleets = append(fleets, Fleet{Name: f})
	}

	return fleets, nil
}

func getFleet(name string) (*Fleet, bool){
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