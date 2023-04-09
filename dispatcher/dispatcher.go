package dispatcher

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
)

type Dispatcher struct {
	ListeningAddress string
	ListeningPort    string

	cmp Camp
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
	d.cmp.Leader.Name = "leader"
	d.cmp.Status = STATUS_STOPPED
	d.cmp.Leader.AuthenticationHash = HashOf("password")
	d.cmp.Soldiers = make([]Soldier, 0)
}
