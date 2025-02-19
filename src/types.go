package main

import (
	"os"
	"path"
)

type Revision struct {
	Name  		string
	Body  		[]byte
	FleetName 	string
}

type Fleet struct {
	Name     		string
}

type Machine struct {
	SerialNumber 	string
	FleetName    	string
	RevisionName    string
	LastSync		string
}

func (r *Revision) save() error {
	filename := path.Join("fleets", r.FleetName, r.Name+".yml")
	return os.WriteFile(filename, r.Body, 0600)
}

func (r *Revision) Fleet() *Fleet {
	f, _ := getFleet(r.FleetName)

	return f
}

func (r *Revision) Machines() []*Machine {
	return getMachinesOnRevision(r.FleetName, r.Name)
}

func (f *Fleet) RevisionNames() []string {
	rs, _ := getRevisionNames(f.Name)

	return rs
}

func (f *Fleet) CurrentRevision() *Revision {
	r, err := getRevision(f.Name, f.CurrentRevisionName())

	if err != nil {

	}

	return r
}

func (f *Fleet) CurrentRevisionName() string {
	return getLatestRevisionName(f.Name)
}

func (f *Fleet) NextRevisionName() string {
	return getNextRevisionName(f.Name)
}

func (f *Fleet) Machines() []*Machine {
	return getMachinesInFleet(f.Name)
}

func (m *Machine) Revision() *Revision {
	r, err := getRevision(m.FleetName, m.RevisionName)

	if err != nil {

	}

	return r
}