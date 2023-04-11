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

func (camp CampAPI) Equals(c2 CampAPI) bool {
	if camp.Leader.Name != c2.Leader.Name {
		return false
	}
	if len(camp.Soldiers) != len(c2.Soldiers) {
		return false
	}
	for i, soldier := range camp.Soldiers {
		if soldier.Name != c2.Soldiers[i].Name {
			return false
		}
		if soldier.Ip != c2.Soldiers[i].Ip {
			return false
		}
	}
	if camp.Settings.Status != c2.Settings.Status {
		return false
	}
	if camp.Settings.VictimServer != c2.Settings.VictimServer {
		return false
	}
	if camp.Settings.DDOSType != c2.Settings.DDOSType {
		return false
	}
	return true
}
