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
		if m.FleetName == fleet || (fleet == defaultFleet && m.FleetName == ""){
			ms = append(ms, machines[i])
		}
	}

	return ms
}

func getMachinesOnRevision(fleet string, revision string) []*Machine {
	ms := make([]*Machine, 0)
	for i, m := range machines {
		if m.RevisionName == revision && m.FleetName == fleet {
			ms = append(ms, machines[i])
		}
	}

	return ms
}