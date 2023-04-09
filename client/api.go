package client

type LeaderAPI struct {
	Name string `json:"name"`
}

type SoldierAPI struct {
	Name string `json:"name"`
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
