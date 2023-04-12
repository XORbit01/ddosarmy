package client

type LeaderAPI struct {
	Name string `json:"name"`
}

type SoldierAPI struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type CampSettings struct {
	Status       string `json:"status"`
	VictimServer string `json:"victim"`
	DDOSType     string `json:"ddos_type"`
}

type CampAPI struct {
	Leader   LeaderAPI    `json:"leader"`
	Soldiers []SoldierAPI `json:"soldiers"`
	Settings CampSettings `json:"camp_settings"`
}

func compareMembers(first []SoldierAPI, second []SoldierAPI) ([]SoldierAPI, []SoldierAPI) {
	firstMap := make(map[string]string)
	secondMap := make(map[string]string)

	// Count members in first slice
	for _, member := range first {
		firstMap[member.Ip] = member.Name
	}

	// Count members in second slice
	for _, member := range second {
		secondMap[member.Ip] = member.Name
	}

	added := make([]SoldierAPI, 0)
	removed := make([]SoldierAPI, 0)

	// Check for removed members
	for ip, name := range firstMap {
		if _, ok := secondMap[ip]; !ok {
			removed = append(removed, SoldierAPI{Name: name, Ip: ip})
		}
	}

	// Check for added members
	for ip, name := range secondMap {
		if _, ok := firstMap[ip]; !ok {
			added = append(added, SoldierAPI{Name: name, Ip: ip})
		}
	}

	return added, removed
}

func (camp CampAPI) Equals(c2 CampAPI) (yes bool, message string) {
	message = ""
	yes = true
	if camp.Leader.Name != c2.Leader.Name {
		yes = false
		message += "leader changed: " + c2.Leader.Name + "\n"
	}
	// check if set of soldiers is the same if not give me the difference
	added, removed := compareMembers(c2.Soldiers, camp.Soldiers)
	if len(added) > 0 {
		yes = false
		message += "soldier joined: "
		for _, soldier := range added {
			message += soldier.Name + "\n"
		}
	}
	if len(removed) > 0 {
		yes = false
		message += "soldier left: "
		for _, soldier := range removed {
			message += soldier.Name + "\n"
		}
	}

	if camp.Settings.Status != c2.Settings.Status {
		yes = false
		message += "+" + camp.Settings.Status + "\n"
	}
	if camp.Settings.VictimServer != c2.Settings.VictimServer {
		yes = false
		message += "victim changed: " + camp.Settings.VictimServer + "\n"
	}
	if camp.Settings.DDOSType != c2.Settings.DDOSType {
		yes = false
		message += "ddos type changed: " + camp.Settings.DDOSType + "\n"
	}
	return yes, message
}
