package main

func getMachine(serial string) (m *Machine, ok bool) {
	for i, m := range machines {
		if m.SerialNumber == serial {
			return machines[i], true
		}
	}

	return nil, false
}

func getMachinesInFleet(fleet string) []*Machine {

	ms := make([]*Machine, 0)
	for i, m := range machines {
		if m.fleetName == fleet || (fleet == defaultFleet && m.fleetName == ""){
			ms = append(ms, machines[i])
		}
	}

	return ms
}

func getMachinesOnRevision(fleet string, revision string) []*Machine {
	ms := make([]*Machine, 0)
	for i, m := range machines {
		if m.RevisionName == revision && (m.fleetName == fleet || (m.fleetName == "" && fleet == "default")) {
			ms = append(ms, machines[i])
		}
	}

	return ms
}