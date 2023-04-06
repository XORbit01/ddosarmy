package camp

import (
	"log"
	"sync"
	"time"
)

type Client struct {
	Name string `json:"name"`
}

type Leader struct {
	Client
}

type Soldier struct {
	Client
	Ip                   string
	LastOrderRequestTime time.Time
}

var instance *Camp
var once sync.Once

type JsonCamp struct {
	Leader   Leader    `json:"leader"`
	Soldiers []Soldier `json:"soldiers"`
}

const NOTHING = "nothing"
const ATTACK = "attack"
const STOP = "stop"

type Camp struct {
	Leader       Leader
	Soldiers     []Soldier
	VictimServer string
	Status       string
}

func NewCamp(leaderName string, victimServer string) *Camp {
	once.Do(func() {
		instance = &Camp{
			Leader:       Leader{Client{Name: leaderName}},
			VictimServer: victimServer,
			Status:       NOTHING,
		}
	})
	return instance
}

func GetCamp() *Camp {
	if instance == nil {
		c := NewCamp("default", "http://localhost:8080")
		return c
	}
	return instance
}

func (c *Camp) SetLeader(l Leader) {
	c.Leader = l
}

func (c *Camp) AddSoldier(s Soldier) {
	c.Soldiers = append(c.Soldiers, s)
}

func (c *Camp) RemoveSoldier(name string) bool {
	for i, sl := range c.Soldiers {
		if sl.Name == name {
			c.Soldiers = append(c.Soldiers[:i], c.Soldiers[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Camp) IsSoldierInCamp(name string) bool {
	for _, sl := range c.Soldiers {
		if sl.Name == name {
			return true
		}
	}
	return false
}

func (c *Camp) Equals(other Camp) bool {
	if c.VictimServer != other.VictimServer {
		return false
	}
	if c.Leader.Name != other.Leader.Name {
		return false
	}
	if len(c.Soldiers) != len(other.Soldiers) {
		return false
	}
	if c.Status != other.Status {
		return false
	}
	for i, sl := range c.Soldiers {
		if sl.Name != other.Soldiers[i].Name {
			return false
		}
	}
	return true
}

func (c *Camp) GetSoldier(name string) *Soldier {
	for i, sl := range c.Soldiers {
		if sl.Name == name {
			return &c.Soldiers[i]
		}
	}
	return nil
}

func (c *Camp) ScanAndRemoveTimeOutedSoldiers() {
	for _, sl := range c.Soldiers {
		if sl.isTimeOutedSoldier() {
			c.RemoveSoldier(sl.Name)
			log.Printf("Soldier %s is timeouted and removed from camp", sl.Name)
		}
	}
}

func (s *Soldier) isTimeOutedSoldier() bool {
	if s.LastOrderRequestTime.IsZero() {
		s.LastOrderRequestTime = time.Now()
	}
	return time.Now().Sub(s.LastOrderRequestTime) > 5*time.Second
}

func (s *Soldier) UpdateLastOrderRequestTime() {
	(*s).LastOrderRequestTime = time.Now()
}

func (s *Soldier) GetLastOrderRequestTime() time.Time {
	return s.LastOrderRequestTime
}
